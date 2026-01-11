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
	_, err = db.Exec(query, 9, "유니버셜 웨이트 3", 78)

	query2 := `
		INSERT INTO tarot_cards (deck_idx, seq, card_name, is_major, img_path, created_at) VALUES
(9, 1, 'THE MAGICIAN', 1, '/cards/uni3/1.jpg',strftime('%s', 'now') ),
(9, 2, 'THE HIGH PRIESTESS', 1, '/cards/uni3/2.jpg',strftime('%s', 'now') ),
(9, 3, 'THE EMPRESS', 1, '/cards/uni3/3.jpg',strftime('%s', 'now') ),
(9, 4, 'THE EMPEROR', 1, '/cards/uni3/4.jpg',strftime('%s', 'now') ),
(9, 5, 'THE HIEROPHANT', 1, '/cards/uni3/5.jpg',strftime('%s', 'now') ),
(9, 6, 'THE LOVERS', 1, '/cards/uni3/6.jpg',strftime('%s', 'now') ),
(9, 7, 'THE CHARIOT', 1, '/cards/uni3/7.jpg',strftime('%s', 'now') ),
(9, 8, 'STRENGTH', 1, '/cards/uni3/8.jpg',strftime('%s', 'now') ),
(9, 9, 'THE HERMIT', 1, '/cards/uni3/9.jpg',strftime('%s', 'now') ),
(9, 10, 'WHEEL OF FORTUNE', 1, '/cards/uni3/10.jpg',strftime('%s', 'now') ),
(9, 11, 'JUSTICE', 1, '/cards/uni3/11.jpg',strftime('%s', 'now') ),
(9, 12, 'THE HANGED MAN', 1, '/cards/uni3/12.jpg',strftime('%s', 'now') ),
(9, 13, 'DEATH', 1, '/cards/uni3/13.jpg',strftime('%s', 'now') ),
(9, 14, 'TEMPERANCE', 1, '/cards/uni3/14.jpg',strftime('%s', 'now') ),
(9, 15, 'THE DEVIL', 1, '/cards/uni3/15.jpg',strftime('%s', 'now') ),
(9, 16, 'THE TOWER', 1, '/cards/uni3/16.jpg',strftime('%s', 'now') ),
(9, 17, 'THE STAR', 1, '/cards/uni3/17.jpg',strftime('%s', 'now') ),
(9, 18, 'THE MOON', 1, '/cards/uni3/18.jpg',strftime('%s', 'now') ),
(9, 19, 'THE SUN', 1, '/cards/uni3/19.jpg',strftime('%s', 'now') ),
(9, 20, 'JUDGEMENT', 1, '/cards/uni3/20.jpg',strftime('%s', 'now') ),
(9, 0, 'THE FOOL', 1, '/cards/uni3/0.jpg',strftime('%s', 'now') ),
(9, 21, 'THE WORLD', 1, '/cards/uni3/21.jpg',strftime('%s', 'now') ),
(9, 22, 'Ace of Wands', 0, '/cards/uni3/22.jpg',strftime('%s', 'now') ),
(9, 23, '2 of Wands', 0, '/cards/uni3/23.jpg',strftime('%s', 'now') ),
(9, 24, '3 of Wands', 0, '/cards/uni3/24.jpg',strftime('%s', 'now') ),
(9, 25, '4 of Wands', 0, '/cards/uni3/25.jpg',strftime('%s', 'now') ),
(9, 26, '5 of Wands', 0, '/cards/uni3/26.jpg',strftime('%s', 'now') ),
(9, 27, '6 of Wands', 0, '/cards/uni3/27.jpg',strftime('%s', 'now') ),
(9, 28, '7 of Wands', 0, '/cards/uni3/28.jpg',strftime('%s', 'now') ),
(9, 29, '8 of Wands', 0, '/cards/uni3/29.jpg',strftime('%s', 'now') ),
(9, 30, '9 of Wands', 0, '/cards/uni3/30.jpg',strftime('%s', 'now') ),
(9, 31, '10 of Wands', 0, '/cards/uni3/31.jpg',strftime('%s', 'now') ),
(9, 32, 'Page of Wands', 0, '/cards/uni3/32.jpg',strftime('%s', 'now') ),
(9, 33, 'Knight of Wands', 0, '/cards/uni3/33.jpg',strftime('%s', 'now') ),
(9, 34, 'Queen of Wands', 0, '/cards/uni3/34.jpg',strftime('%s', 'now') ),
(9, 35, 'King of Wands', 0, '/cards/uni3/35.jpg',strftime('%s', 'now') ),
(9, 36, 'Ace of Cups', 0, '/cards/uni3/36.jpg',strftime('%s', 'now') ),
(9, 37, '2 of Cups', 0, '/cards/uni3/37.jpg',strftime('%s', 'now') ),
(9, 38, '3 of Cups', 0, '/cards/uni3/38.jpg',strftime('%s', 'now') ),
(9, 39, '4 of Cups', 0, '/cards/uni3/39.jpg',strftime('%s', 'now') ),
(9, 40, '5 of Cups', 0, '/cards/uni3/40.jpg',strftime('%s', 'now') ),
(9, 41, '6 of Cups', 0, '/cards/uni3/41.jpg',strftime('%s', 'now') ),
(9, 42, '7 of Cups', 0, '/cards/uni3/42.jpg',strftime('%s', 'now') ),
(9, 43, '8 of Cups', 0, '/cards/uni3/43.jpg',strftime('%s', 'now') ),
(9, 44, '9 of Cups', 0, '/cards/uni3/44.jpg',strftime('%s', 'now') ),
(9, 45, '10 of Cups', 0, '/cards/uni3/45.jpg',strftime('%s', 'now') ),
(9, 46, 'Page of Cups', 0, '/cards/uni3/46.jpg',strftime('%s', 'now') ),
(9, 47, 'Knight of Cups', 0, '/cards/uni3/47.jpg',strftime('%s', 'now') ),
(9, 48, 'Queen of Cups', 0, '/cards/uni3/48.jpg',strftime('%s', 'now') ),
(9, 49, 'King of Cups', 0, '/cards/uni3/49.jpg',strftime('%s', 'now') ),
(9, 50, 'Ace of Swords', 0, '/cards/uni3/50.jpg',strftime('%s', 'now') ),
(9, 51, '2 of Swords', 0, '/cards/uni3/51.jpg',strftime('%s', 'now') ),
(9, 52, '3 of Swords', 0, '/cards/uni3/52.jpg',strftime('%s', 'now') ),
(9, 53, '4 of Swords', 0, '/cards/uni3/53.jpg',strftime('%s', 'now') ),
(9, 54, '5 of Swords', 0, '/cards/uni3/54.jpg',strftime('%s', 'now') ),
(9, 55, '6 of Swords', 0, '/cards/uni3/55.jpg',strftime('%s', 'now') ),
(9, 56, '7 of Swords', 0, '/cards/uni3/56.jpg',strftime('%s', 'now') ),
(9, 57, '8 of Swords', 0, '/cards/uni3/57.jpg',strftime('%s', 'now') ),
(9, 58, '9 of Swords', 0, '/cards/uni3/58.jpg',strftime('%s', 'now') ),
(9, 59, '10 of Swords', 0, '/cards/uni3/59.jpg',strftime('%s', 'now') ),
(9, 60, 'Page of Swords', 0, '/cards/uni3/60.jpg',strftime('%s', 'now') ),
(9, 61, 'Knight of Swords', 0, '/cards/uni3/61.jpg',strftime('%s', 'now') ),
(9, 62, 'Queen of Swords', 0, '/cards/uni3/62.jpg',strftime('%s', 'now') ),
(9, 63, 'King of Swords', 0, '/cards/uni3/63.jpg',strftime('%s', 'now') ),
(9, 64, 'Ace of Pentacles', 0, '/cards/uni3/64.jpg',strftime('%s', 'now') ),
(9, 65, '2 of Pentacles', 0, '/cards/uni3/65.jpg',strftime('%s', 'now') ),
(9, 66, '3 of Pentacles', 0, '/cards/uni3/66.jpg',strftime('%s', 'now') ),
(9, 67, '4 of Pentacles', 0, '/cards/uni3/67.jpg',strftime('%s', 'now') ),
(9, 68, '5 of Pentacles', 0, '/cards/uni3/68.jpg',strftime('%s', 'now') ),
(9, 69, '6 of Pentacles', 0, '/cards/uni3/69.jpg',strftime('%s', 'now') ),
(9, 70, '7 of Pentacles', 0, '/cards/uni3/70.jpg',strftime('%s', 'now') ),
(9, 71, '8 of Pentacles', 0, '/cards/uni3/71.jpg',strftime('%s', 'now') ),
(9, 72, '9 of Pentacles', 0, '/cards/uni3/72.jpg',strftime('%s', 'now') ),
(9, 73, '10 of Pentacles', 0, '/cards/uni3/73.jpg',strftime('%s', 'now') ),
(9, 74, 'Page of Pentacles', 0, '/cards/uni3/74.jpg',strftime('%s', 'now') ),
(9, 75, 'Knight of Pentacles', 0, '/cards/uni3/75.jpg',strftime('%s', 'now') ),
(9, 76, 'Queen of Pentacles', 0, '/cards/uni3/76.jpg',strftime('%s', 'now') ),
(9, 77, 'King of Pentacles', 0, '/cards/uni3/77.jpg',strftime('%s', 'now') )
	`
	_, err = db.Exec(query2)

	fmt.Println("🚀 마이그레이션 완료!")
}
