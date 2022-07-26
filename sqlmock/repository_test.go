package sqlmock_test

import (
	"database/sql"
	"regexp"
	"time"

	"dbtest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ = Describe("Repository", func() {
	var repository *dbtest.Repository
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var db *sql.DB
		var err error

		// db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)) // use equal matcher
		db, mock, err = sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())

		gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		Expect(err).ShouldNot(HaveOccurred())

		repository = &dbtest.Repository{Db: gdb}
	})
	AfterEach(func() {
		err := mock.ExpectationsWereMet() // make sure all expectations were met
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("list all", func() {
		It("empty", func() {
			const sqlSelectAll = `SELECT * FROM "blogs"`
			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
				WillReturnRows(sqlmock.NewRows(nil))

			l, err := repository.ListAll()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(l).Should(BeEmpty())
		})

	})

	Context("load", func() {
		It("found", func() {
			blog := &dbtest.Blog{
				ID:        1,
				Title:     "post",
				Content:   "hello",
				Tags:      pq.StringArray{"go", "golang"},
				CreatedAt: time.Now(),
			}

			rows := sqlmock.
				NewRows([]string{"id", "title", "content", "tags", "created_at"}).
				AddRow(blog.ID, blog.Title, blog.Content, blog.Tags, blog.CreatedAt)

			const sqlSelectOne = `SELECT * FROM "blogs" WHERE id = $1 ORDER BY "blogs"."id" LIMIT 1`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(blog.ID).
				WillReturnRows(rows)

			dbBlog, err := repository.Load(blog.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dbBlog).Should(Equal(blog))
		})

		It("not found", func() {
			// ignore sql match
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := repository.Load(1)
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

	Context("list", func() {
		It("found", func() {
			rows := sqlmock.
				NewRows([]string{"id", "title", "content", "tags", "created_at"}).
				AddRow(1, "post 1", "hello 1", nil, time.Now()).
				AddRow(2, "post 2", "hello 2", pq.StringArray{"go"}, time.Now())

			// limit/offset is not parameter
			const sqlSelectFirstTen = `SELECT * FROM "blogs" LIMIT 10`
			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectFirstTen)).WillReturnRows(rows)

			l, err := repository.List(0, 10)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(l).Should(HaveLen(2))
			Expect(l[0].Tags).Should(BeEmpty())
			Expect(l[1].Tags).Should(Equal(pq.StringArray{"go"}))
			Expect(l[1].ID).Should(BeEquivalentTo(2)) // use BeEquivalentTo
		})
		It("not found", func() {
			// ignore sql match
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			l, err := repository.List(0, 10)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(l).Should(BeEmpty())
		})
	})

	Context("save", func() {
		var blog *dbtest.Blog
		BeforeEach(func() {
			blog = &dbtest.Blog{
				Title:     "post",
				Content:   "hello",
				Tags:      pq.StringArray{"a", "b"},
				CreatedAt: time.Now(),
			}
		})

		It("update", func() {
			const sqlUpdate = `
					UPDATE "blogs" SET "title"=$1,"content"=$2,"tags"=$3,"created_at"=$4 WHERE "id" = $5`
			const sqlSelectOne = `
					SELECT * FROM "blogs" WHERE "blogs"."id" = $1 ORDER BY "blogs"."id" ASC LIMIT 1`

			blog.ID = 1
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
				WithArgs(blog.Title, blog.Content, blog.Tags, blog.CreatedAt, blog.ID).
				WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()

			err := repository.Save(blog)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("insert", func() {
			// gorm use query instead of exec
			// https://github.com/DATA-DOG/go-sqlmock/issues/118
			const sqlInsert = `
					INSERT INTO "blogs" ("title","content","tags","created_at") 
						VALUES ($1,$2,$3,$4) RETURNING "id"`
			const newId = 1
			mock.ExpectBegin() // start transaction
			mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
				WithArgs(blog.Title, blog.Content, blog.Tags, blog.CreatedAt).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
			mock.ExpectCommit() // commit transaction

			Expect(blog.ID).Should(BeZero())

			err := repository.Save(blog)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(blog.ID).Should(BeEquivalentTo(newId))
		})

	})

	Context("search by title", func() {
		It("found", func() {
			rows := sqlmock.
				NewRows([]string{"id", "title", "content", "tags", "created_at"}).
				AddRow(1, "post 1", "hello 1", nil, time.Now())

			// limit/offset is not parameter
			const sqlSearch = `
				SELECT * FROM "blogs" 
				WHERE title like $1 
				LIMIT 10`
			const q = "os"

			mock.ExpectQuery(regexp.QuoteMeta(sqlSearch)).
				WithArgs("%" + q + "%").
				WillReturnRows(rows)

			l, err := repository.SearchByTitle(q, 0, 10)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(l).Should(HaveLen(1))
			Expect(l[0].Title).Should(ContainSubstring(q))
		})
	})
})
