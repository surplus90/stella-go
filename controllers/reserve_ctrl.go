package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"stella-tarot/database" // go.mod의 모듈명 확인 필요
	"stella-tarot/models"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// 화면에 보여주기 위한 확장 구조체
type ReservationView struct {
	models.Reservation
	DeckName        string
	ReservationDate string
}
type DisplayCard struct {
	models.TarotCard
	Position int // 뽑은 순서 (1번, 2번...)
}

// 보안 키 생성 함수 (함수 내부에 있거나 별도로 정의되어야 함)
func generateEncKey() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// 1. 예약 폼 페이지 렌더링
func RenderReservePage(c echo.Context) error {
	rows, err := database.DB.Query("SELECT idx, deck_name, amount_cards FROM tarot_decks")
	if err != nil {
		return c.String(http.StatusInternalServerError, "덱 목록 조회 실패")
	}
	defer rows.Close()

	// 1. 현재 시간 구하기
	now := time.Now()
	defaultTime := now.Add(time.Hour * 1)
	formattedTime := defaultTime.Format("2006-01-02T15:04")

	var decks []models.TarotDeck
	for rows.Next() {
		var d models.TarotDeck
		// CreatedAt은 쿼리에서 뺐습니다 (화면 표시에 불필요)
		rows.Scan(&d.Idx, &d.DeckName, &d.AmountCards)
		decks = append(decks, d)
	}

	return c.Render(http.StatusOK, "reserve.html", map[string]interface{}{
		"IsAdmin":     true,
		"Decks":       decks,
		"DefaultTime": formattedTime,
	})
}

// 2. 예약 정보 저장
func SaveReservation(c echo.Context) error {
	userName := c.FormValue("user_name")
	deckIdx, _ := strconv.ParseInt(c.FormValue("deck_idx"), 10, 64)
	amountCards, _ := strconv.Atoi(c.FormValue("amount_cards"))
	selectedCards, _ := strconv.Atoi(c.FormValue("selected_cards"))
	wayToArray, _ := strconv.Atoi(c.FormValue("way_to_array"))
	dateStr := c.FormValue("reservation_at")

	// 날짜 변환
	loc, _ := time.LoadLocation("Asia/Seoul")
	t, _ := time.ParseInLocation("2006-01-02T15:04", dateStr, loc)
	reservationAt := t.Unix()

	encKey := generateEncKey()
	now := time.Now().Unix()

	_, err := database.DB.Exec(`
		INSERT INTO reservation (
			user_name, deck_idx, amount_cards, selected_cards, 
			way_to_array, reservation_at, enc_key, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userName, deckIdx, amountCards, selectedCards,
		wayToArray, reservationAt, encKey, now,
	)

	if err != nil {
		return c.String(http.StatusInternalServerError, "DB 저장 실패")
	}

	return c.Redirect(http.StatusSeeOther, "/reserve-list")
}

// 3. 예약 목록 화면 랜더링
func RenderReserveListPage(c echo.Context) error {
	// 1. 페이지 번호 가져오기 (기본값 1)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	// 2. 전체 데이터 개수 조회 (페이지네이션 UI용)
	var totalCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM reservation").Scan(&totalCount)
	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	// 3. 20개씩 끊어서 조회 (LEFT JOIN으로 데이터 누락 방지)
	// Scan 에러 방지를 위해 COALESCE(null 처리) 추가
	query := `
        SELECT 
			r.idx, r.user_name, r.amount_cards, r.selected_cards, 
			r.way_to_array, r.reservation_at, r.enc_key, r.created_at,
            d.deck_name
        FROM reservation r
        LEFT JOIN tarot_decks d ON r.deck_idx = d.idx
        ORDER BY r.idx DESC
        LIMIT ? OFFSET ?
    `
	rows, err := database.DB.Query(query, pageSize, offset)
	if err != nil {
		return c.String(http.StatusInternalServerError, "목록 조회 실패: "+err.Error())
	}
	defer rows.Close()

	var reservations []ReservationView
	for rows.Next() {
		var rv ReservationView
		err := rows.Scan(
			&rv.Idx, &rv.UserName, &rv.AmountCards, &rv.SelectedCards,
			&rv.WayToArray, &rv.ReservationAt, &rv.EncKey, &rv.CreatedAt,
			&rv.DeckName,
		)
		if err != nil {
			continue
		}

		// 2. 시간 포맷 변환 (Unix -> "2006-01-02 15:04")
		tm := time.Unix(rv.ReservationAt, 0)
		rv.ReservationDate = tm.Format("2006-01-02 15:04")

		reservations = append(reservations, rv)
	}

	return c.Render(http.StatusOK, "reserve-list.html", map[string]interface{}{
		"IsAdmin":      true,
		"Reservations": reservations,
		"CurrentPage":  page,
		"TotalPage":    totalPage, // totalPage 사용 확인
		"TotalCount":   totalCount,
		"HasPrev":      page > 1,
		"HasNext":      page < totalPage,
		"PrevPage":     page - 1,
		"NextPage":     page + 1,
	})
}

// 4. 예약 상세 정보 및 뽑은 카드 확인
func RenderReserveDetailPage(c echo.Context) error {
	idx := c.Param("idx")

	// 1. 예약 및 덱 정보 가져오기
	query := `
		SELECT 
			r.idx, r.user_name, r.amount_cards, r.selected_cards, 
			r.way_to_array, r.reservation_at, r.enc_key, r.created_at,
			r.deck_idx, d.deck_name
		FROM reservation r
		LEFT JOIN tarot_decks d ON r.deck_idx = d.idx
		WHERE r.idx = ?
	`
	var rv ReservationView
	err := database.DB.QueryRow(query, idx).Scan(
		&rv.Idx, &rv.UserName, &rv.AmountCards, &rv.SelectedCards,
		&rv.WayToArray, &rv.ReservationAt, &rv.EncKey, &rv.CreatedAt,
		&rv.DeckIdx, &rv.DeckName,
	)
	if err != nil {
		return c.String(http.StatusNotFound, "예약을 찾을 수 없습니다.")
	}

	// 시간 포맷팅
	rv.ReservationDate = time.Unix(rv.ReservationAt, 0).Format("2006-01-02 15:04")

	// 2. 고객이 뽑은 카드 번호 가져오기
	var pickedCards string
	err = database.DB.QueryRow("SELECT cards FROM pick_cards WHERE reservation_idx = ?", idx).Scan(&pickedCards)

	var displayCards []DisplayCard
	if err == nil && pickedCards != "" {
		// 여기서 cardIds 변수를 선언합니다! (문자열 "1,5,12" -> ["1", "5", "12"])
		cardIds := strings.Split(pickedCards, ",")

		for i, idStr := range cardIds {
			var card DisplayCard
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}

			// string을 int로 변환
			seqInt, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}

			// TarotCard 테이블에서 정보 조회
			cardQuery := `SELECT card_name, img_path FROM tarot_cards WHERE deck_idx = ? AND seq = ?`
			err = database.DB.QueryRow(cardQuery, rv.DeckIdx, seqInt).Scan(&card.CardName, &card.ImgPath)

			if err != nil {
				fmt.Printf("[조회실패] Deck:%d, Seq:%d, Err:%v\n", rv.DeckIdx, seqInt, err)
				continue
			}

			card.Position = i + 1
			displayCards = append(displayCards, card)
		}
	}

	return c.Render(http.StatusOK, "reserve-detail.html", map[string]interface{}{
		"IsAdmin":      true,
		"Res":          rv,
		"DisplayCards": displayCards, // 선택된 카드 번호 리스트
	})
}

// 5. 예약 삭제 처리
func DeleteReservation(c echo.Context) error {
	idx := c.Param("idx")

	// DB에서 삭제 쿼리 실행
	_, err := database.DB.Exec("DELETE FROM reservation WHERE idx = ?", idx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "삭제 실패: "+err.Error())
	}

	// 삭제 후 예약 목록 페이지로 리다이렉트
	return c.Redirect(http.StatusSeeOther, "/reserve-list")
}

// 6. 카드 추가
func AddCardCount(c echo.Context) error {
	idx := c.Param("idx")
	addCountStr := c.FormValue("add_count")

	// 1. 입력값 숫자로 변환
	addCount, err := strconv.Atoi(addCountStr)
	if err != nil || addCount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "잘못된 숫자입니다."})
	}

	// 2. DB 업데이트: 기존 값에 더하기 (SQL의 + 연산 활용)
	query := `UPDATE reservation SET selected_cards = selected_cards + ? WHERE idx = ?`
	_, err = database.DB.Exec(query, addCount, idx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "업데이트 실패"})
	}

	// 3. 다시 상세 페이지로 리다이렉트
	return c.Redirect(http.StatusSeeOther, "/reserve-detail/"+idx)
}
