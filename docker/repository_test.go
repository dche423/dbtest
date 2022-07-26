package docker_test

import (
	"time"

	"dbtest"

	"github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var repo *dbtest.Repository
	BeforeEach(func() {
		repo = &dbtest.Repository{Db: Db}
		err := repo.Migrate() // auto create tables
		Ω(err).To(Succeed())

		sampleData := &dbtest.Blog{
			Title:     "post",
			Content:   "hello",
			Tags:      []string{"a", "b"},
			CreatedAt: time.Now(),
		}
		// create sample data
		err = Db.Create(sampleData).Error
		Ω(err).To(Succeed())
	})
	Context("Load", func() {
		It("Found", func() {
			blog, err := repo.Load(1)

			Ω(err).To(Succeed())
			Ω(blog.Content).To(Equal("hello"))
			Ω(blog.Tags).To(Equal(pq.StringArray{"a", "b"}))
		})
		It("Not Found", func() {
			_, err := repo.Load(999)
			Ω(err).To(HaveOccurred())
		})
	})
	It("ListAll", func() {
		l, err := repo.ListAll()
		Ω(err).To(Succeed())
		Ω(l).To(HaveLen(1))
	})
	It("List", func() {
		l, err := repo.List(0, 10)
		Ω(err).To(Succeed())
		Ω(l).To(HaveLen(1))
	})
	Context("Save", func() {
		It("Create", func() {
			blog := &dbtest.Blog{
				Title:     "post2",
				Content:   "hello",
				Tags:      []string{"a", "b"},
				CreatedAt: time.Now(),
			}
			err := repo.Save(blog)
			Ω(err).To(Succeed())
			Ω(blog.ID).To(BeEquivalentTo(2))
		})
		It("Update", func() {
			blog, err := repo.Load(1)
			Ω(err).To(Succeed())

			blog.Title = "foo"
			err = repo.Save(blog)
			Ω(err).To(Succeed())
		})
	})
	It("Delete", func() {
		err := repo.Delete(1)
		Ω(err).To(Succeed())
		_, err = repo.Load(1)
		Ω(err).To(HaveOccurred())
	})
	DescribeTable("SearchByTitle",
		func(q string, found int) {
			l, err := repo.SearchByTitle(q, 0, 10)
			Ω(err).To(Succeed())
			Ω(l).To(HaveLen(found))
		},
		Entry("found", "post", 1),
		Entry("not found", "bar", 0),
	)
})
