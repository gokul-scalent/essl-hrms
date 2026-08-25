
package service

import (
	"context"
	"strconv"
	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type EmployeeServiceImpl struct {
	employeeRepo EmployeeRepo
}

func NewEmployeeServiceImpl(employeeRepo EmployeeRepo) (*EmployeeServiceImpl, error) {
	return &EmployeeServiceImpl{
    employeeRepo: employeeRepo,
    }, nil
}


func (s *EmployeeServiceImpl) CreateEmployee(ctx context.Context, employee entity.Employee) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>employee: create employee started", reqID)

	employeeID,errResp := s.employeeRepo.CreateEmployee(ctx, employee)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0,errResp
	}

	log.Info("core>service>employee: create employee completed & employee id is "+strconv.Itoa(employeeID), reqID)
	return employeeID,nil
}



func (s *EmployeeServiceImpl) PartialUpdateEmployee(ctx context.Context, employee entity.Employee) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>employee: partila update employee started for employee id "+strconv.Itoa(employee.ID), reqID)

	errResp := s.employeeRepo.PartialUpdateEmployee(ctx, employee)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>employee: update employee completed for employee id "+strconv.Itoa(employee.ID), reqID)
	return nil
}

	

func (s *EmployeeServiceImpl) UpdateEmployee(ctx context.Context, employee entity.Employee) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>employee: update employee started for employee id "+strconv.Itoa(employee.ID), reqID)

	errResp := s.employeeRepo.UpdateEmployee(ctx, employee)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>employee: update employee completed for employee id "+strconv.Itoa(employee.ID), reqID)
	return nil
}

	

func (s *EmployeeServiceImpl) DeleteEmployee(ctx context.Context, employeeID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>employee: delete employee started for employee id "+strconv.Itoa(employeeID), reqID)

	
	errResp := s.employeeRepo.DeleteEmployee(ctx, employeeID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>employee: delete employee completed for employee id "+strconv.Itoa(employeeID), reqID)
	return nil
}

	

func (s *EmployeeServiceImpl) GetEmployeebyID(ctx context.Context, employeeID int) (entity.Employee, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>employee: get employee started for employee id "+strconv.Itoa(employeeID), reqID)

	employee, errResp := s.employeeRepo.GetEmployeebyID(ctx, employeeID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return entity.Employee{}, errResp
	}

	
	log.Info("core>service>employee: employee fetched successfully for employee id "+strconv.Itoa(employeeID), reqID)
	return employee, nil
}

	

func (s *EmployeeServiceImpl) ListEmployee(ctx context.Context, filter *filters.ListFilter) (int, []entity.Employee, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>employee: employee list started", reqID)

	totalRecords, employeesEntity, errResp := s.employeeRepo.ListEmployee(ctx, filter)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, nil, errResp
	}

	log.Info("core>service>employee: employee list completed", reqID)
	return totalRecords, employeesEntity, nil
}

