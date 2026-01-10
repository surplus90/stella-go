package controllers

import (
	"math/rand"
	"net/http"
	"stella-tarot/database"
	// "stella-tarot/models"
	"time"
	"strings"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RenderSelectCardsPage(c echo.Context) error {
	encKey := c.Param("enc_key")

	// 1. 예약 정보와 덱 정보 JOIN 조회 (DeckIdx 추가)
	var data struct {
		Idx           int64
		UserName      string
		DisplayCount  int    // 화면에 깔아줄 장수
		PickCount     int    // 설정된 총 선택 장수 (Reservation.SelectedCards)
		TotalDeck     int    // 덱 전체 장수
		EncKey        string
		DeckIdx       int    // 카드 조회를 위한 덱 번호
	}

	query := `
		SELECT 
			r.idx, r.user_name, r.amount_cards, r.selected_cards, r.enc_key,
			d.amount_cards, r.deck_idx
		FROM reservation r
		JOIN tarot_decks d ON r.deck_idx = d.idx
		WHERE r.enc_key = ?
	`
	err := database.DB.QueryRow(query, encKey).Scan(
		&data.Idx, &data.UserName, &data.DisplayCount, &data.PickCount, &data.EncKey,
		&data.TotalDeck, &data.DeckIdx,
	)
	if err != nil {
		return c.String(http.StatusNotFound, "예약 정보를 찾을 수 없습니다.")
	}

	// 2. 이미 선택한 카드가 있는지 확인
	var pickedCardsStr string
	// 에러 무시 (데이터가 없을 수 있음)
	database.DB.QueryRow("SELECT cards FROM pick_cards WHERE reservation_idx = ?", data.Idx).Scan(&pickedCardsStr)

	var pickedSeqs []string
	if pickedCardsStr != "" {
		pickedSeqs = strings.Split(pickedCardsStr, ",")
	}
	pickedCount := len(pickedSeqs)

	// --- [분기점 1] 설정된 장수만큼 다 뽑았을 경우: 결과 화면 노출 ---
	if pickedCount >= data.PickCount && pickedCount > 0 {
		var displayCards []DisplayCard
		for i, seqStr := range pickedSeqs {
			var card DisplayCard
			cardQuery := `SELECT card_name, img_path FROM tarot_cards WHERE deck_idx = ? AND seq = ?`
			err := database.DB.QueryRow(cardQuery, data.DeckIdx, strings.TrimSpace(seqStr)).Scan(&card.CardName, &card.ImgPath)
			if err == nil {
				card.Position = i + 1
				displayCards = append(displayCards, card)
			}
		}

		return c.Render(http.StatusOK, "select-cards.html", map[string]interface{}{
			"HideNav": true,
			"IsResult":     true, // 템플릿에서 결과 화면임을 판별하는 플래그
			"UserName":     data.UserName,
			"DisplayCards": displayCards,
			"PickCount":    data.PickCount,
		})
	}

	// --- [분기점 2] 아직 뽑아야 할 카드가 남았을 경우: 셔플 및 선택 화면 노출 ---

	// 3. 이미 뽑은 카드 리스트를 Map으로 변환 (제외 대상 필터링용)
	pickedMap := make(map[int]bool)
	for _, s := range pickedSeqs {
		val, _ := strconv.Atoi(strings.TrimSpace(s))
		pickedMap[val] = true
	}

	// 4. 전체 덱 중 이미 뽑은 카드를 제외하고 생성
	var availablePool []int
	for i := 0; i < data.TotalDeck; i++ {
		if !pickedMap[i] { // 이미 뽑은 Seq(i)가 아닐 때만 풀에 추가
			availablePool = append(availablePool, i)
		}
	}

	// 5. 남은 카드 풀 셔플
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(availablePool), func(i, j int) {
		availablePool[i], availablePool[j] = availablePool[j], availablePool[i]
	})

	// 6. 화면에 노출할 카드 추출 (이미 뽑은 카드를 제외한 나머지가 DisplayCount보다 작을 수 있음 대응)
	displayLimit := data.DisplayCount
	if len(availablePool) < displayLimit {
		displayLimit = len(availablePool)
	}
	finalDisplayCards := availablePool[:displayLimit]

	// 추가로 뽑아야 할 장수 계산
	needToPick := data.PickCount - pickedCount

	return c.Render(http.StatusOK, "select-cards.html", map[string]interface{}{
		"HideNav": true,
		"IsResult":     false,
		"Idx":          data.Idx,
		"UserName":     data.UserName,
		"DisplayCards": finalDisplayCards, // 셔플된 Seq 리스트
		"PickCount":    needToPick,        // 추가로 뽑을 장수
		"AlreadyPick":  pickedCount,       // 기존에 뽑은 장수
		"EncKey":       data.EncKey,
	})
}

// SubmitSelectedCards - 선택한 카드 저장
func SubmitSelectedCardsOri(c echo.Context) error {
	encKey := c.Param("enc_key")
	resIdx := c.FormValue("reservation_idx")
	selectedValues := c.FormValue("selected_values") // JS에서 콤마로 합쳐서 보낼 예정

	if selectedValues == "" {
		return c.String(http.StatusBadRequest, "카드를 선택해주세요.")
	}

	// PickCards 테이블에 저장
	_, err := database.DB.Exec(
		"INSERT INTO pick_cards (reservation_idx, enc_key, cards, created_at) VALUES (?, ?, ?, ?)",
		resIdx, encKey, selectedValues, time.Now().Unix(),
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, "저장 실패: "+err.Error())
	}

	return c.HTML(http.StatusOK, "<h2>카드가 성공적으로 제출되었습니다. 상담사를 기다려주세요!</h2>")
}

func SubmitSelectedCards(c echo.Context) error {
	encKey := c.Param("enc_key")
	resIdx := c.FormValue("reservation_idx")
	selectedValues := c.FormValue("selected_values") // 예: "10,22"

	if selectedValues == "" {
		return c.String(http.StatusBadRequest, "카드를 선택해주세요.")
	}

	// 1. 기존에 저장된 카드가 있는지 확인
	var existingCards string
	err := database.DB.QueryRow("SELECT cards FROM pick_cards WHERE reservation_idx = ?", resIdx).Scan(&existingCards)

	if err == nil {
		// [CASE 1] 기존 데이터가 있는 경우 -> 기존 카드 + 콤마 + 신규 카드 (UPDATE)
		newCombinedCards := existingCards + "," + selectedValues
		_, err = database.DB.Exec(
			"UPDATE pick_cards SET cards = ?, created_at = ? WHERE reservation_idx = ?",
			newCombinedCards, time.Now().Unix(), resIdx,
		)
		if err != nil {
			return c.String(http.StatusInternalServerError, "업데이트 실패: "+err.Error())
		}
	} else {
		// [CASE 2] 기존 데이터가 없는 경우 (처음 뽑는 경우) -> 신규 생성 (INSERT)
		_, err = database.DB.Exec(
			"INSERT INTO pick_cards (reservation_idx, enc_key, cards, created_at) VALUES (?, ?, ?, ?)",
			resIdx, encKey, selectedValues, time.Now().Unix(),
		)
		if err != nil {
			return c.String(http.StatusInternalServerError, "저장 실패: "+err.Error())
		}
	}

	// 성공 후 안내 메시지 (고객이 결과 화면을 바로 볼 수 있도록 리다이렉트하는 것도 좋습니다)
	return c.Redirect(http.StatusSeeOther, "/select-cards/"+encKey)
}