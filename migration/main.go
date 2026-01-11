package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	// 2. SQLite 연결 (파일 위치 주의: migration 폴더 안에서 실행하므로 ../stella.db)
	db, err := sql.Open("sqlite", "../sqlite/stella.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// SQLite 마이그레이션 쿼리 (Unix Timestamp 정수형 저장)
	query := `
		INSERT INTO tarot_decks (idx, deck_name, amount_cards, created_at)
		VALUES (?, ?, ?, strftime('%s', 'now'))
	`
	_, err = db.Exec(query, 8, "I CHING of LOVE", 64)

	query2 := `
		INSERT INTO tarot_cards (deck_idx, seq, card_name, is_major, img_path, created_at) VALUES
(8,0,1,1,'/cards/ichingoflove/0.jpg',strftime('%s', 'now') ),
(8,1,2,1,'/cards/ichingoflove/1.jpg',strftime('%s', 'now') ),
(8,2,3,1,'/cards/ichingoflove/2.jpg',strftime('%s', 'now') ),
(8,3,4,1,'/cards/ichingoflove/3.jpg',strftime('%s', 'now') ),
(8,4,5,1,'/cards/ichingoflove/4.jpg',strftime('%s', 'now') ),
(8,5,6,1,'/cards/ichingoflove/5.jpg',strftime('%s', 'now') ),
(8,6,7,1,'/cards/ichingoflove/6.jpg',strftime('%s', 'now') ),
(8,7,8,1,'/cards/ichingoflove/7.jpg',strftime('%s', 'now') ),
(8,8,9,1,'/cards/ichingoflove/8.jpg',strftime('%s', 'now') ),
(8,9,10,1,'/cards/ichingoflove/9.jpg',strftime('%s', 'now') ),
(8,10,11,1,'/cards/ichingoflove/10.jpg',strftime('%s', 'now') ),
(8,11,12,1,'/cards/ichingoflove/11.jpg',strftime('%s', 'now') ),
(8,12,13,1,'/cards/ichingoflove/12.jpg',strftime('%s', 'now') ),
(8,13,14,1,'/cards/ichingoflove/13.jpg',strftime('%s', 'now') ),
(8,14,15,1,'/cards/ichingoflove/14.jpg',strftime('%s', 'now') ),
(8,15,16,1,'/cards/ichingoflove/15.jpg',strftime('%s', 'now') ),
(8,16,17,1,'/cards/ichingoflove/16.jpg',strftime('%s', 'now') ),
(8,17,18,1,'/cards/ichingoflove/17.jpg',strftime('%s', 'now') ),
(8,18,19,1,'/cards/ichingoflove/18.jpg',strftime('%s', 'now') ),
(8,19,20,1,'/cards/ichingoflove/19.jpg',strftime('%s', 'now') ),
(8,20,21,1,'/cards/ichingoflove/20.jpg',strftime('%s', 'now') ),
(8,21,22,1,'/cards/ichingoflove/21.jpg',strftime('%s', 'now') ),
(8,22,23,1,'/cards/ichingoflove/22.jpg',strftime('%s', 'now') ),
(8,23,24,1,'/cards/ichingoflove/23.jpg',strftime('%s', 'now') ),
(8,24,25,1,'/cards/ichingoflove/24.jpg',strftime('%s', 'now') ),
(8,25,26,1,'/cards/ichingoflove/25.jpg',strftime('%s', 'now') ),
(8,26,27,1,'/cards/ichingoflove/26.jpg',strftime('%s', 'now') ),
(8,27,28,1,'/cards/ichingoflove/27.jpg',strftime('%s', 'now') ),
(8,28,29,1,'/cards/ichingoflove/28.jpg',strftime('%s', 'now') ),
(8,29,30,1,'/cards/ichingoflove/29.jpg',strftime('%s', 'now') ),
(8,30,31,1,'/cards/ichingoflove/30.jpg',strftime('%s', 'now') ),
(8,31,32,1,'/cards/ichingoflove/31.jpg',strftime('%s', 'now') ),
(8,32,33,1,'/cards/ichingoflove/32.jpg',strftime('%s', 'now') ),
(8,33,34,1,'/cards/ichingoflove/33.jpg',strftime('%s', 'now') ),
(8,34,35,1,'/cards/ichingoflove/34.jpg',strftime('%s', 'now') ),
(8,35,36,1,'/cards/ichingoflove/35.jpg',strftime('%s', 'now') ),
(8,36,37,1,'/cards/ichingoflove/36.jpg',strftime('%s', 'now') ),
(8,37,38,1,'/cards/ichingoflove/37.jpg',strftime('%s', 'now') ),
(8,38,39,1,'/cards/ichingoflove/38.jpg',strftime('%s', 'now') ),
(8,39,40,1,'/cards/ichingoflove/39.jpg',strftime('%s', 'now') ),
(8,40,41,1,'/cards/ichingoflove/40.jpg',strftime('%s', 'now') ),
(8,41,42,1,'/cards/ichingoflove/41.jpg',strftime('%s', 'now') ),
(8,42,43,1,'/cards/ichingoflove/42.jpg',strftime('%s', 'now') ),
(8,43,44,1,'/cards/ichingoflove/43.jpg',strftime('%s', 'now') ),
(8,44,45,1,'/cards/ichingoflove/44.jpg',strftime('%s', 'now') ),
(8,45,46,1,'/cards/ichingoflove/45.jpg',strftime('%s', 'now') ),
(8,46,47,1,'/cards/ichingoflove/46.jpg',strftime('%s', 'now') ),
(8,47,48,1,'/cards/ichingoflove/47.jpg',strftime('%s', 'now') ),
(8,48,49,1,'/cards/ichingoflove/48.jpg',strftime('%s', 'now') ),
(8,49,50,1,'/cards/ichingoflove/49.jpg',strftime('%s', 'now') ),
(8,50,51,1,'/cards/ichingoflove/50.jpg',strftime('%s', 'now') ),
(8,51,52,1,'/cards/ichingoflove/51.jpg',strftime('%s', 'now') ),
(8,52,53,1,'/cards/ichingoflove/52.jpg',strftime('%s', 'now') ),
(8,53,54,1,'/cards/ichingoflove/53.jpg',strftime('%s', 'now') ),
(8,54,55,1,'/cards/ichingoflove/54.jpg',strftime('%s', 'now') ),
(8,55,56,1,'/cards/ichingoflove/55.jpg',strftime('%s', 'now') ),
(8,56,57,1,'/cards/ichingoflove/56.jpg',strftime('%s', 'now') ),
(8,57,58,1,'/cards/ichingoflove/57.jpg',strftime('%s', 'now') ),
(8,58,59,1,'/cards/ichingoflove/58.jpg',strftime('%s', 'now') ),
(8,59,60,1,'/cards/ichingoflove/59.jpg',strftime('%s', 'now') ),
(8,60,61,1,'/cards/ichingoflove/60.jpg',strftime('%s', 'now') ),
(8,61,62,1,'/cards/ichingoflove/61.jpg',strftime('%s', 'now') ),
(8,62,63,1,'/cards/ichingoflove/62.jpg',strftime('%s', 'now') ),
(8,63,64,1,'/cards/ichingoflove/63.jpg',strftime('%s', 'now') )
	`
	_, err = db.Exec(query2)

	fmt.Println("🚀 마이그레이션 완료!")
}
