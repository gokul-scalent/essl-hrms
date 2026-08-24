package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	"github.com/scalent.io/scalent-hrms/pkg/reacher"
	"github.com/scalent.io/scalent-hrms/pkg/validation"
)

type LeadServiceImpl struct {
	leadRepo        LeadRepo
	emailListRepo   EmailListRepo
	reacherClient   *reacher.Client
	UserSettingRepo UserSettingRepo
}

func NewLeadServiceImpl(leadRepo LeadRepo, emailListRepo EmailListRepo, reacherClient *reacher.Client, UserSettingRepo UserSettingRepo) (*LeadServiceImpl, error) {
	return &LeadServiceImpl{
		leadRepo:        leadRepo,
		emailListRepo:   emailListRepo,
		reacherClient:   reacherClient,
		UserSettingRepo: UserSettingRepo,
	}, nil
}

func (s *LeadServiceImpl) CreateLead(ctx context.Context, lead entity.Lead, req *coreAPIModel.CreateLeadRequest) (int, errors.Response) {

	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: create lead started", reqID)

	session, err := mailoraContext.GetSessionFromContext(ctx)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseUnauthorizedError(errors.SESSION_ERROR)
	}

	userID := session.UserID

	// Validate Email List
	emailList, errResp := s.emailListRepo.GetEmailListbyID(ctx, lead.EmailList.ID, userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	switch req.Type {

	case commonConstants.UPLOAD_TYPE_SINGLE:
		// duplicate check fun call
		exists, errResp := s.leadRepo.IsLeadExists(
			ctx,
			lead.EmailList.ID,
			lead.Email,
		)
		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return 0, errResp
		}

		if exists {
			return 0, errors.ResponseBadRequestError("lead already exists")
		}
		lead.Priority = emailList.Priority
		lead.VerificationStatus = commonConstants.VERIFICATION_STATUS_PENDING
		lead.IsSafe = commonConstants.LEAD_STATUS_UNKNOWN
		lead.FinalStatus = commonConstants.LEAD_STATUS_UNKNOWN
		lead.IsReachable = commonConstants.LEAD_STATUS_UNKNOWN

		leadID, errResp := s.leadRepo.CreateLead(ctx, lead)
		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return 0, errResp
		}

		errResp = s.leadRepo.UpdateEmailListCounts(ctx, lead.EmailList.ID)
		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return 0, errResp
		}

		log.Info("core>service>lead: create lead completed & lead id is "+strconv.Itoa(leadID), reqID)
		return leadID, nil

	case commonConstants.UPLOAD_TYPE_CSV:

		file, err := req.File.Open()
		if err != nil {
			log.Error(err.Error(), reqID)
			return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
		defer file.Close()

		reader := csv.NewReader(file)

		// Read header row
		headers, err := reader.Read()
		if err != nil {
			log.Error(err.Error(), reqID)
			return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}

		// Supported CSV header aliases
		fieldMap := map[string]string{
			// Email
			"email":         "Email",
			"email address": "Email",
			"e-mail":        "Email",
			"to email":      "Email",

			// First Name
			"first name": "FirstName",
			"firstname":  "FirstName",
			"name":       "FirstName",
			"given name": "FirstName",
			"to name":    "FirstName",
			// Last Name
			"last name":   "LastName",
			"lastname":    "LastName",
			"surname":     "LastName",
			"family name": "LastName",

			// Company
			"company":      "Company",
			"company name": "Company",
			"organization": "Company",
			"organisation": "Company",
			"employer":     "Company",
			"to company":   "Company",

			// Job Title
			"job title":   "JobTitle",
			"title":       "JobTitle",
			"designation": "JobTitle",
			"position":    "JobTitle",

			// City
			"city": "City",
			"town": "City",

			// Country
			"country": "Country",
			"nation":  "Country",

			// Industry
			"industry": "Industry",
			"sector":   "Industry",
		}

		// Build: LeadField -> CSV Column Index
		columnIndex := make(map[string]int)

		for i, h := range headers {
			header := strings.ToLower(strings.TrimSpace(h))

			if field, ok := fieldMap[header]; ok {
				columnIndex[field] = i
			}
		}

		// Email is mandatory
		if _, ok := columnIndex["Email"]; !ok {
			return 0, errors.ResponseBadRequestError("CSV must contain an Email column")
		}

		// Generic getter
		get := func(row []string, field string) string {
			if idx, ok := columnIndex[field]; ok {
				if idx < len(row) {
					return strings.TrimSpace(row[idx])
				}
			}
			return ""
		}

		insertedCount := 0
		lastLeadID := 0

		for {
			row, err := reader.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Error(err.Error(), reqID)
				return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
			}

			if len(row) == 0 {
				continue
			}

			email := get(row, "Email")

			if email == "" {
				continue
			}

			csvLead := entity.Lead{
				EmailList: entity.EmailList{
					ID: lead.EmailList.ID,
				},

				Email:              email,
				EmailProvIDer:      validation.GetEmailProviderFromMX(email),
				FirstName:          get(row, "FirstName"),
				LastName:           get(row, "LastName"),
				JobTitle:           get(row, "JobTitle"),
				Company:            get(row, "Company"),
				City:               get(row, "City"),
				Country:            get(row, "Country"),
				Industry:           get(row, "Industry"),
				Priority:           emailList.Priority,
				VerificationStatus: commonConstants.VERIFICATION_STATUS_PENDING,
				IsSafe:             commonConstants.LEAD_STATUS_UNKNOWN,
				FinalStatus:        commonConstants.LEAD_STATUS_UNKNOWN,
				IsReachable:        commonConstants.LEAD_STATUS_UNKNOWN,
			}

			// Duplicate check
			exists, errResp := s.leadRepo.IsLeadExists(
				ctx,
				csvLead.EmailList.ID,
				csvLead.Email,
			)

			if errResp != nil {
				log.Error(errResp.Error(), reqID)
				continue
			}

			if exists {
				continue
			}

			leadID, errResp := s.leadRepo.CreateLead(ctx, csvLead)

			if errResp != nil {
				log.Error(errResp.Error(), reqID)
				continue
			}

			insertedCount++
			lastLeadID = leadID
		}

		if insertedCount == 0 {
			return 0, errors.ResponseBadRequestError("no valid emails found")
		}

		errResp = s.leadRepo.UpdateEmailListCounts(ctx, lead.EmailList.ID)
		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return 0, errResp
		}

		log.Info("core>service>lead: csv upload completed", reqID)

		return lastLeadID, nil

	default:
		return 0, errors.ResponseBadRequestError("invalid upload type")
	}
}

func (s *LeadServiceImpl) PartialUpdateLead(ctx context.Context, lead entity.Lead) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: partial update lead started for lead id "+strconv.Itoa(lead.ID), reqID)

	// Get existing lead before update
	existingLead, errResp := s.leadRepo.GetLeadbyID(ctx, lead.ID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	// Check whether email changed
	emailChanged := lead.Email != "" &&
		!strings.EqualFold(
			strings.TrimSpace(existingLead.Email),
			strings.TrimSpace(lead.Email),
		)

	// Update lead
	errResp = s.leadRepo.PartialUpdateLead(ctx, lead, emailChanged)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	// If email changed, recalculate parent email list counts
	if emailChanged {
		errResp = s.leadRepo.UpdateEmailListCounts(
			ctx,
			existingLead.EmailList.ID,
		)
		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return errResp
		}
	}

	log.Info("core>service>lead: update lead completed for lead id "+strconv.Itoa(lead.ID), reqID)
	return nil
}

func (s *LeadServiceImpl) UpdateLead(ctx context.Context, lead entity.Lead) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: update lead started for lead id "+strconv.Itoa(lead.ID), reqID)

	errResp := s.leadRepo.UpdateLead(ctx, lead)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>lead: update lead completed for lead id "+strconv.Itoa(lead.ID), reqID)
	return nil
}

func (s *LeadServiceImpl) DeleteLead(ctx context.Context, leadID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: delete lead started for lead id "+strconv.Itoa(leadID), reqID)

	errResp := s.leadRepo.DeleteLead(ctx, leadID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>lead: delete lead completed for lead id "+strconv.Itoa(leadID), reqID)
	return nil
}

func (s *LeadServiceImpl) GetLeadbyID(ctx context.Context, leadID int) (entity.Lead, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: get lead started for lead id "+strconv.Itoa(leadID), reqID)

	lead, errResp := s.leadRepo.GetLeadbyID(ctx, leadID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return entity.Lead{}, errResp
	}

	log.Info("core>service>lead: lead fetched successfully for lead id "+strconv.Itoa(leadID), reqID)
	return lead, nil
}

func (s *LeadServiceImpl) ListLead(ctx context.Context, emailListID int, filter *filters.ListFilter) (int, []entity.Lead, entity.LeadStatusCount, errors.Response) {

	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: lead list started", reqID)

	totalRecords, leads, statusCount, errResp := s.leadRepo.ListLead(ctx, emailListID, filter)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, nil, entity.LeadStatusCount{}, errResp
	}

	log.Info("core>service>lead: lead list completed", reqID)
	return totalRecords, leads, statusCount, nil
}

// get the safe leads for download purpose
func (s *LeadServiceImpl) GetSafeLeads(ctx context.Context, emailListID int, status string) ([]entity.Lead, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: get safe leads started", reqID)

	leads, errResp := s.leadRepo.GetSafeLeads(ctx, emailListID, status)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return nil, errResp
	}

	log.Info("core>service>lead: get safe leads completed", reqID)
	return leads, nil
}

func (s *LeadServiceImpl) VerifyPendingLeads(ctx context.Context) (time.Duration, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: VerifyPendingLeads started", reqID)

	// Get the highest priority email list having pending leads
	emailList, errResp := s.emailListRepo.GetNextEmailListForVerification(ctx)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	// No pending email list
	if emailList == nil {
		log.Info("No email list found for verification", reqID)
		return 0, nil
	}

	// Get verification interval from DB
	verificationInterval, errResp :=
		s.UserSettingRepo.GetVerificationInterval(ctx, emailList.User.ID)

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	// Convert seconds to time.Duration
	interval := time.Duration(verificationInterval) * time.Second

	log.Info(fmt.Sprintf("Verification interval for user %d is %v", emailList.User.ID, interval), reqID)
	log.Info(fmt.Sprintf("Email list ID: %d", emailList.ID), reqID)
	// Fetch only one pending lead from that email list
	leads, errResp := s.leadRepo.GetPendingLeadsByEmailListID(
		ctx,
		emailList.ID,
		1,
	)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	// when timeout records has time for reverification or Nothing pending
	if len(leads) == 0 {
		log.Info("No pending leads ready for verification", reqID)
		return interval, nil
	}

	lead := leads[0]
	log.Info("Verifying email : "+lead.Email, reqID)

	// Call Reacher
	resp, err := s.reacherClient.VerifyEmail(ctx, lead.Email)
	if err != nil {
		log.Error(err.Error(), reqID)
		if strings.Contains(strings.ToLower(err.Error()), "timeout") ||
			strings.Contains(strings.ToLower(err.Error()), "deadline exceeded") {
			resp = &reacher.VerifyEmailResponse{
				IsReachable: "TIMEOUT",
			}

			errResp = s.leadRepo.UpdateLeadVerification(ctx, lead.ID, resp)
			if errResp != nil {
				log.Error(errResp.Error(), reqID)
				return 0, errResp
			}

			errResp = s.leadRepo.UpdateEmailListCounts(ctx, emailList.ID)
			if errResp != nil {
				log.Error(errResp.Error(), reqID)
				return 0, errResp
			}

			return interval, nil
		}
		return 0, errors.ResponseInternalServerError(err.Error())
	}

	// Update verification result
	errResp = s.leadRepo.UpdateLeadVerification(
		ctx,
		lead.ID,
		resp,
	)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	// Update counts of parent email list
	errResp = s.leadRepo.UpdateEmailListCounts(
		ctx,
		emailList.ID,
	)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	log.Info("Verified email : "+lead.Email, reqID)
	log.Info("core>service>lead: VerifyPendingLeads completed", reqID)

	return interval, nil
}

func (s *LeadServiceImpl) HasPendingVerification(ctx context.Context) (bool, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: HasPendingVerification started", reqID)

	emailList, errResp := s.emailListRepo.GetNextEmailListForVerification(ctx)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return false, errResp
	}

	hasPending := emailList != nil
	log.Info(fmt.Sprintf("Has pending verification Status: %v", hasPending), reqID)
	log.Info("core>service>lead: HasPendingVerification completed", reqID)
	return hasPending, nil
}

// re verify the lead
func (s *LeadServiceImpl) ReverifyLead(ctx context.Context, leadID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>lead: ReverifyLead started for lead id "+strconv.Itoa(leadID), reqID)

	// Get parent email list ID before re-verification
	emailListID, errResp := s.leadRepo.GetLeadEmailListID(ctx, leadID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	// Reset lead verification state
	errResp = s.leadRepo.ReverifyLead(ctx, leadID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	// Recalculate parent email list counts
	errResp = s.leadRepo.UpdateEmailListCounts(ctx, emailListID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>lead: ReverifyLead completed for lead id "+strconv.Itoa(leadID), reqID)
	return nil
}
