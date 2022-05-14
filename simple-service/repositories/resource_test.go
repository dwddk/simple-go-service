package repositories_test

import (
	"context"
	"errors"
	"regexp"

	"github.com/addme96/simple-go-service/simple-service/entities"
	"github.com/addme96/simple-go-service/simple-service/repositories"
	"github.com/addme96/simple-go-service/simple-service/repositories/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pashagolub/pgxmock"
)

var _ = Describe("Resource", func() {
	var (
		ctrl     *gomock.Controller
		mockDB   *mocks.MockDB
		repo     *repositories.Resource
		ctx      context.Context
		mockConn pgxmock.PgxConnIface
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockDB(ctrl)
		repo = repositories.NewResource(mockDB)
		ctx = context.Background()
		mockConn, _ = pgxmock.NewConn()

	})

	Context("Create", func() {
		query := "INSERT into resources (name) VALUES ($1)"

		Context("happy path", func() {
			It("creates the resource", func() {
				By("arranging")
				resourceToCreate := entities.Resource{ID: 101, Name: "Resource Name"}
				mockDB.EXPECT().GetConn(ctx).Times(1).Return(mockConn, nil)
				mockConn.ExpectPrepare("createResource", regexp.QuoteMeta(query)).
					ExpectExec().WithArgs(resourceToCreate.Name).WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mockConn.ExpectClose()

				By("acting")
				err := repo.Create(ctx, resourceToCreate)

				By("asserting")
				Expect(err).NotTo(HaveOccurred())
				Expect(mockConn.ExpectationsWereMet()).To(Succeed())
			})
		})

		Context("not so happy path", func() {
			When("GetConn fails", func() {
				It("returns error", func() {
					By("arranging")
					expectedResource := entities.Resource{ID: 101, Name: "Resource Name"}
					expectedErr := errors.New("some error")
					mockDB.EXPECT().GetConn(ctx).Times(1).Return(nil, expectedErr)

					By("acting")
					err := repo.Create(ctx, expectedResource)

					By("asserting")
					Expect(err).To(Equal(expectedErr))
					Expect(mockConn.ExpectationsWereMet()).To(Succeed())
				})
			})

			When("Prepare fails", func() {
				It("returns error", func() {
					By("arranging")
					mockDB.EXPECT().GetConn(ctx).Times(1).Return(mockConn, nil)
					expectedErr := errors.New("prepare error")
					mockConn.ExpectPrepare("createResource", regexp.QuoteMeta(query)).WillReturnError(expectedErr)
					mockConn.ExpectClose()

					By("acting")
					err := repo.Create(ctx, entities.Resource{})

					By("asserting")
					Expect(err).To(Equal(expectedErr))
					Expect(mockConn.ExpectationsWereMet()).To(Succeed())
				})
			})

			When("Exec fails", func() {
				It("returns error", func() {
					By("arranging")
					expectedResource := entities.Resource{ID: 101, Name: "Resource Name"}
					mockDB.EXPECT().GetConn(ctx).Times(1).Return(mockConn, nil)
					expectedErr := errors.New("exec error")
					mockConn.ExpectPrepare("createResource", regexp.QuoteMeta(query)).
						ExpectExec().WithArgs(expectedResource.Name).WillReturnError(expectedErr)
					mockConn.ExpectClose()

					By("acting")
					err := repo.Create(ctx, expectedResource)

					By("asserting")
					Expect(err).To(Equal(expectedErr))
					Expect(mockConn.ExpectationsWereMet()).To(Succeed())
				})
			})
		})
	})

	Context("Read", func() {
		query := "SELECT id, name FROM resources WHERE id=$1"

		Context("happy path", func() {
			It("reads the resource", func() {
				By("arranging")
				expectedResource := entities.Resource{ID: 101, Name: "Resource Name"}
				mockDB.EXPECT().GetConn(ctx).Times(1).Return(mockConn, nil)
				rows := pgxmock.NewRows([]string{"id", "name"}).AddRow(expectedResource.ID, expectedResource.Name)
				mockConn.ExpectPrepare("readResource", regexp.QuoteMeta(query)).ExpectQuery().
					WithArgs(expectedResource.ID).WillReturnRows(rows)
				mockConn.ExpectClose()

				By("acting")
				res, err := repo.Read(ctx, expectedResource.ID)

				By("asserting")
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(&expectedResource))
				Expect(mockConn.ExpectationsWereMet()).To(Succeed())
			})
		})

		Context("not so happy path", func() {
			When("GetConn fails", func() {
				It("returns error", func() {
					By("arranging")
					expectedErr := errors.New("some error")
					mockDB.EXPECT().GetConn(ctx).Times(1).Return(nil, expectedErr)

					By("acting")
					res, err := repo.Read(ctx, 101)

					By("asserting")
					Expect(err).To(Equal(expectedErr))
					Expect(res).To(BeNil())
					Expect(mockConn.ExpectationsWereMet()).To(Succeed())
				})
			})

			When("Prepare fails", func() {
				It("returns error", func() {
					By("arranging")
					mockDB.EXPECT().GetConn(ctx).Times(1).Return(mockConn, nil)
					expectedErr := errors.New("prepare error")
					mockConn.ExpectPrepare("readResource", regexp.QuoteMeta(query)).WillReturnError(expectedErr)
					mockConn.ExpectClose()

					By("acting")
					res, err := repo.Read(ctx, 101)

					By("asserting")
					Expect(err).To(Equal(expectedErr))
					Expect(res).To(BeNil())
					Expect(mockConn.ExpectationsWereMet()).To(Succeed())
				})
			})

			When("QueryRow fails", func() {
				It("returns error", func() {
					By("arranging")
					resourceID := 101
					mockDB.EXPECT().GetConn(ctx).Times(1).Return(mockConn, nil)
					expectedErr := errors.New("query row error")
					mockConn.ExpectPrepare("readResource", regexp.QuoteMeta(query)).
						ExpectQuery().WithArgs(resourceID).WillReturnError(expectedErr)
					mockConn.ExpectClose()

					By("acting")
					res, err := repo.Read(ctx, resourceID)

					By("asserting")
					Expect(err).To(Equal(expectedErr))
					Expect(res).To(BeNil())
					Expect(mockConn.ExpectationsWereMet()).To(Succeed())
				})
			})
		})
	})

	Context("ReadAll", func() {
		Context("happy path", func() {

		})

		Context("not so happy path", func() {

		})
	})

	Context("Update", func() {
		Context("happy path", func() {

		})

		Context("not so happy path", func() {

		})
	})

	Context("Delete", func() {
		Context("happy path", func() {

		})

		Context("not so happy path", func() {

		})
	})
})