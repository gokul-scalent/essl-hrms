package repo

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scalent.io/scalent-hrms/entity"
	commonConstants "github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	"github.com/scalent.io/scalent-hrms/internal/converter"
	"github.com/scalent.io/scalent-hrms/model"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	filterPkg "github.com/scalent.io/scalent-hrms/pkg/filter"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type EmployeeRepoImpl struct {
	db *sqlx.DB
}

func NewEmployeeRepoImpl(db *sqlx.DB) (*EmployeeRepoImpl, error) {
	return &EmployeeRepoImpl{
		db: db,
	}, nil
}

func (r *EmployeeRepoImpl) CreateEmployee(ctx context.Context, employee entity.Employee) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>employee: CreateEmployee started", reqID)

	query := "INSERT INTO employees (uid, emp_id, emp_name, privilege, password, group_id, card) VALUES(?, ?, ?, ?, ?, ?, ?)"

	result, err := r.db.Exec(query, employee.UID, employee.EmpID, employee.EmpName, employee.Privilege, employee.Password, employee.GroupID, employee.Card)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	employeeID, err := result.LastInsertId()
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>employee: CreateEmployee completed & employee id is "+strconv.Itoa(int(employeeID)), reqID)
	return int(employeeID), nil
}

func (r *EmployeeRepoImpl) PartialUpdateEmployee(ctx context.Context, employee entity.Employee) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>employee:  PartialUpdateEmployee started for employee id "+strconv.Itoa(employee.ID), reqID)

	columns := []string{}
	args := []interface{}{}

	if employee.Card != "" {
		columns = append(columns, "card=?")
		args = append(args, employee.Card)
	}

	if employee.EmpID != "" {
		columns = append(columns, "emp_id=?")
		args = append(args, employee.EmpID)
	}

	if employee.EmpName != "" {
		columns = append(columns, "emp_name=?")
		args = append(args, employee.EmpName)
	}

	if employee.GroupID != "" {
		columns = append(columns, "group_id=?")
		args = append(args, employee.GroupID)
	}

	if employee.Password != "" {
		columns = append(columns, "password=?")
		args = append(args, employee.Password)
	}

	if employee.Privilege != 0 {
		columns = append(columns, "privilege=?")
		args = append(args, employee.Privilege)
	}

	if employee.UID != 0 {
		columns = append(columns, "uid=?")
		args = append(args, employee.UID)
	}

	args = append(args, employee.ID)

	columnStr := strings.Join(columns, ", ")

	if columnStr != "" {
		query := "UPDATE employees SET " + columnStr + " WHERE  id=?  "

		_, err := r.db.Exec(query, args...)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
	}

	log.Info("core>repo>employee: PartialUpdateEmployee completed for employee id "+strconv.Itoa(employee.ID), reqID)
	return nil
}

func (r *EmployeeRepoImpl) UpdateEmployee(ctx context.Context, employee entity.Employee) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>employee: UpdateEmployee started for employee id "+strconv.Itoa(employee.ID), reqID)

	query := "UPDATE employees SET uid=?, emp_id=?, emp_name=?, privilege=?, password=?, group_id=?, card=? WHERE id=?  "

	_, err := r.db.Exec(query, employee.UID, employee.EmpID, employee.EmpName, employee.Privilege, employee.Password, employee.GroupID, employee.Card, employee.ID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>employee: UpdateEmployee completed for employee id "+strconv.Itoa(employee.ID), reqID)
	return nil
}

func (r *EmployeeRepoImpl) DeleteEmployee(ctx context.Context, employeeID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>employee: DeleteEmployee started for employee id "+strconv.Itoa(employeeID), reqID)

	query := "UPDATE employees SET deleted_at = ? WHERE id = ?"

	_, err := r.db.Exec(query, time.Now(), employeeID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>employee: DeleteEmployee completed for employee id "+strconv.Itoa(employeeID), reqID)
	return nil
}

func (r *EmployeeRepoImpl) GetEmployeebyID(ctx context.Context, employeeID int) (entity.Employee, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>employee: GetEmployeebyID started for employee id "+strconv.Itoa(employeeID), reqID)

	query := "SELECT * FROM employees WHERE employees.id=? "

	employeeModel := model.Employee{}
	employeeEntity := entity.Employee{}

	err := r.db.Get(&employeeModel, query, employeeID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return employeeEntity, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	employeeEntity = converter.EmployeeModelToEmployeeEntity(employeeModel)

	log.Info("core>repo>employee: GetEmployeebyID completed for employee id "+strconv.Itoa(employeeID), reqID)
	return employeeEntity, nil
}

func (r *EmployeeRepoImpl) ListEmployee(ctx context.Context, filter *filters.ListFilter) (int, []entity.Employee, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>employee: ListEmployee started", reqID)

	queryStatement := "SELECT * FROM employees "

	modelmap := model.EmployeeModelMap

	whereStr, args := filterPkg.CreateFilterStr(filter.Filters, modelmap)

	// Search employee by ID or name
	if filter.SearchString != "" {
		search := "%" + strings.TrimSpace(filter.SearchString) + "%"

		whereStr = append(
			whereStr,
			"(employees.emp_id LIKE ? OR employees.emp_name LIKE ?)",
		)
		args = append(args, search, search)
	}

	// Soft delete
	whereStr = append(whereStr, "employees.deleted_at IS NULL")

	whereString := strings.Join(whereStr, " AND ")
	whereString = "WHERE " + whereString

	queryStatement += whereString

	sortStr := filterPkg.CreateSortStr(filter.SortOption, modelmap)
	queryStatement += sortStr

	var limitQueryStmt string

	emptySortOption := filters.SortOption{}

	totalRecordQueryStatement := "SELECT COUNT(id) as totalRecords FROM (" + queryStatement + ") as result"

	var count int
	err := r.db.Get(&count, totalRecordQueryStatement, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	if filter.Page == 0 && len(filter.Filters) == 0 &&
		filter.SortOption == emptySortOption &&
		filter.SearchString == "" {

		limitQueryStmt = queryStatement

	} else {
		if filter.Page == 0 {
			filter.Page = 1
		}

		offset := commonConstants.NO_OF_RECORDS_PER_PAGE * (filter.Page - 1)

		limitQueryStmt = queryStatement + " LIMIT ?,?"
		args = append(args, offset, commonConstants.NO_OF_RECORDS_PER_PAGE)
	}

	employeesModel := []model.Employee{}

	err = r.db.Select(&employeesModel, limitQueryStmt, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	employeeEntities := []entity.Employee{}

	for _, employeeModel := range employeesModel {
		employeeEntity := converter.EmployeeModelToEmployeeEntity(employeeModel)
		employeeEntities = append(employeeEntities, employeeEntity)
	}

	log.Info("core>repo>employee: ListEmployee completed", reqID)

	return count, employeeEntities, nil
}

func (r *EmployeeRepoImpl) GetEmployeeDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.Employee, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>employee: GetEmployeeDetails started", reqID)

	selectStr := strings.Join(selectColumns, ", ")
	whereStr := strings.Join(whereColumn, " AND ")

	employeeModel := model.Employee{}

	query := "SELECT " + selectStr + " FROM " + table + " WHERE " + whereStr
	err := r.db.Get(&employeeModel, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	employeeEntity := converter.EmployeeModelToEmployeeEntity(employeeModel)

	log.Info("core>repo>employee: GetEmployeeDetails completed", reqID)
	return &employeeEntity, nil
}
