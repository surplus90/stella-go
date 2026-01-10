package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	_ "modernc.org/sqlite"
)

func main() {
	// 1. Redis 연결
	rdb := redis.NewClient(&redis.Options{
		Addr: "158.247.198.238:6379",
		Password: "Gkstndbs!@34",
		DB:   0,
	})
	ctx := context.Background()

	// 2. SQLite 연결 (파일 위치 주의: migration 폴더 안에서 실행하므로 ../stella.db)
	db, err := sql.Open("sqlite", "../sqlite/stella.db") 
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. Redis에서 모든 키 조회
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		log.Fatal(err)
	}

	for _, resIdx := range keys {
		// Redis List 데이터 (뽑은 카드들)
		cardList, err := rdb.LRange(ctx, resIdx, 0, -1).Result()
		if err != nil || len(cardList) == 0 {
			continue
		}
		cardsJoined := strings.Join(cardList, ",")

		// SQLite의 reservation 테이블에서 enc_key 매칭
		var encKey sql.NullString
		_ = db.QueryRow("SELECT enc_key FROM reservation WHERE idx = ?", resIdx).Scan(&encKey)

		finalKey := ""
		if encKey.Valid {
			finalKey = encKey.String
		}

		// SQLite 마이그레이션 쿼리 (Unix Timestamp 정수형 저장)
		query := `
			INSERT INTO pick_cards (reservation_idx, enc_key, cards, created_at)
			VALUES (?, ?, ?, strftime('%s', 'now'))
		`
		_, err = db.Exec(query, resIdx, finalKey, cardsJoined)
		if err != nil {
			log.Printf("❌ 실패 [idx:%s]: %v", resIdx, err)
		} else {
			fmt.Printf("✅ 성공 [idx:%s, key:%s]: %s\n", resIdx, finalKey, cardsJoined)
		}
	}
	fmt.Println("🚀 마이그레이션 완료!")
}