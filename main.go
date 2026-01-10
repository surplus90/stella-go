package main

import (
	"os"
	"log"
	"net/http"
	"html/template"
	"io"
	"stella-tarot/controllers"
	"stella-tarot/database"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

// Echo 템플릿 렌더러 정의
type TemplateRenderer struct {}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    // 1. 함수 정의
    funcMap := template.FuncMap{
        "add": func(a, b int) int {
            return a + b
        },
    }

    // 2. 템플릿 생성 시 .Funcs(funcMap)을 반드시 연결해야 합니다.
    tmpl, err := template.New("layout").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/"+name)
    if err != nil {
        return err
    }

    // 3. 실행
    return tmpl.ExecuteTemplate(w, "layout", data)
}

// 파일에서 비밀번호를 읽어오는 함수
func loadPassword(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("비밀번호 파일을 읽을 수 없습니다: %v", err)
	}
	// 읽어온 내용의 앞뒤 공백(줄바꿈 등)을 제거하고 반환
	return strings.TrimSpace(string(content))
}

func main() {
	e := echo.New()

	// 서버 시작 시 비밀번호 로드
	adminPassword := loadPassword(".env")

	// 정적 파일 설정 (이미지, CSS 등)
	e.Static("/static", "static")
	e.Renderer = &TemplateRenderer{}

	// 1. 세션 설정 (비밀키는 원하는 대로 변경하세요)
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("stk"))))

	// 2. 인증 미들웨어
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()

			// 고객 허용 경로 (메인, 카드선택, 정적파일, 로그인 페이지)
			if path == "/" || strings.HasPrefix(path, "/select-cards") || 
               strings.HasPrefix(path, "/static") || path == "/login" {
				return next(c)
			}

			// 관리자 권한 체크 (세션 확인)
			sess, _ := session.Get("admin-session", c)
			if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
				// 인증되지 않았으면 로그인 페이지로 이동
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			return next(c)
		}
	})

	// 3. 로그인 처리 라우트
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login.html", nil)
	})

	e.POST("/login", func(c echo.Context) error {
		password := c.FormValue("password")

		// 비밀번호 확인 (원하는 비밀번호로 수정하세요)
		if password == adminPassword {
			sess, _ := session.Get("admin-session", c)
			sess.Values["authenticated"] = true
			sess.Save(c.Request(), c.Response())
			return c.Redirect(http.StatusSeeOther, "/reserve-list")
		}

		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "비밀번호가 일치하지 않습니다.",
		})
	})
    
    // 로그아웃
    e.GET("/logout", func(c echo.Context) error {
		sess, _ := session.Get("admin-session", c)
		sess.Values["authenticated"] = false
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, "/")
	})

	// DB 연결
	database.InitDB()

	// --- 라우팅 ---
	e.GET("/", func(c echo.Context) error {
		sess, _ := session.Get("admin-session", c)
		isAdmin, _ := sess.Values["authenticated"].(bool) // 세션에서 인증 여부 확인

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"IsAdmin": isAdmin,
		})
	})

	// 예약 관련
	e.GET("/reserve", controllers.RenderReservePage)
	e.POST("/save-reservation", controllers.SaveReservation)
	e.GET("/reserve-list", controllers.RenderReserveListPage)
	e.GET("/reserve-detail/:idx", controllers.RenderReserveDetailPage)
	e.POST("/delete-reservation/:idx", controllers.DeleteReservation)

	// 카드 뽑기
	e.GET("/select-cards/:enc_key", controllers.RenderSelectCardsPage)
	e.POST("/select-cards/:enc_key", controllers.SubmitSelectedCards)
	e.POST("/add-card-count/:idx", controllers.AddCardCount)

	e.Logger.Fatal(e.Start(":8080"))
}