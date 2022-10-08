package repository

import (
	"context"
	studentDomain "layout_2/internal/domain/student"
	"layout_2/internal/utils"

	"go.uber.org/dig"

	"gorm.io/gorm"
)

type StudentRepositoryCond struct {
	dig.In

	DB *gorm.DB `name:"dbM"`
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(cond StudentRepositoryCond) (studentDomain.StudentRepository, error) {
	result := &studentRepository{
		db: cond.DB,
	}

	if err := result.createTable(); err != nil {
		return nil, err
	}

	return result, nil
}

func (sr *studentRepository) createTable() error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancelCtx()

	if sr.db.WithContext(ctx).Migrator().HasTable(&studentDomain.Student{}) {
		return nil
	}

	err := sr.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='學生表'").
		WithContext(ctx).
		AutoMigrate(&studentDomain.Student{})
	if err != nil {
		return err
	}

	return nil
}

func (sr *studentRepository) Create(student studentDomain.Student) (uint, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancelCtx()

	result := sr.db.WithContext(ctx).Create(&student)
	if result.Error != nil {
		return 0, result.Error
	}

	return student.ID, nil
}

func (sr *studentRepository) UpdateID(cond studentDomain.StudentCond, student studentDomain.Student) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancelCtx()

	result := sr.db.Model(&student).Where("`id` = ?", cond.ID).WithContext(ctx).
		Updates(map[string]interface{}{
			"name":   student.Name,
			"gender": student.Gender,
			"status": student.Status,
		})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (sr *studentRepository) GetID(cond studentDomain.StudentCond) (studentDomain.Student, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancelCtx()

	var student studentDomain.Student

	result := sr.db.WithContext(ctx).First(&student, cond.ID)
	if result.Error != nil {
		return studentDomain.Student{}, result.Error
	}

	return student, nil
}

func (sr *studentRepository) Get(cond studentDomain.StudentCond) ([]studentDomain.Student, int64, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancelCtx()

	var students []studentDomain.Student

	// 注意: 不要覆寫回 sr.db，會紀錄上次 where 條件資料，造成 SQL 錯誤
	db := sr.db

	db = utils.SQLAppend(db, utils.IsNotZero(int(cond.ID)), "`id` = ?", cond.ID)

	db = utils.SQLAppend(db, utils.IsNotEmpty(cond.Name), "`name` like ?", "%"+cond.Name+"%")
	db = utils.SQLAppend(db, utils.IsNotNegativeOne(cond.Gender), "`gender` = ?", cond.Gender)

	db = utils.SQLAppend(db, utils.IsNotNegativeOne(cond.Status), "`status` = ?", cond.Status)

	var count int64
	resultCount := db.Model(&studentDomain.Student{}).WithContext(ctx).Count(&count)
	if resultCount.Error != nil {
		return nil, 0, resultCount.Error
	}

	db = utils.SQLPagination(db, cond.GetRowCount(), cond.GetOffset())

	result := db.WithContext(ctx).Find(&students)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return students, count, nil
}
