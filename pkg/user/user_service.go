package user

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/entities"
	"Go-Starter-Template/internal/utils"
	"Go-Starter-Template/internal/utils/storage"
	jwtService "Go-Starter-Template/pkg/jwt"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, req domain.UserRegisterRequest) (domain.UserRegisterResponse, error)
		Login(ctx context.Context, req domain.UserLoginRequest) (domain.UserLoginResponse, error)
		GetProfile(ctx context.Context, userID string) (domain.UserProfileResponse, error)
		UpdateProfile(ctx context.Context, req domain.UpdateUserRequest, userID string) error
		PostEducation(ctx context.Context, req domain.PostUserEducationRequest, userID string) error
		DeleteEducation(ctx context.Context, educationID string) error
		UpdateEducation(ctx context.Context, req domain.UpdateUserEducationRequest, userID string) error
		PostExperience(ctx context.Context, req domain.PostUserExperienceRequest, userID string) error
		UpdateExperience(ctx context.Context, req domain.UpdateUserExperienceRequest, userID string) error
		DeleteExperience(ctx context.Context, experienceID string) error
		PostSkill(ctx context.Context, req domain.PostUserSkillRequest, userID string) error
		DeleteSkill(ctx context.Context, skillID string) error
		GetSkills(ctx context.Context) ([]domain.SkillsResponse, error)
		SearchUser(ctx context.Context, query domain.UserSearchRequest) ([]domain.UserSearchResponse, error)
	}

	userService struct {
		userRepository UserRepository
		awsS3          storage.AwsS3
		jwtService     jwtService.JWTService
	}
)

func NewUserService(userRepository UserRepository, awsS3 storage.AwsS3, jwtService jwtService.JWTService) UserService {
	return &userService{userRepository: userRepository, awsS3: awsS3, jwtService: jwtService}
}

var VerifyEmailRoute = "api/verify_email/user"

func (s *userService) RegisterUser(ctx context.Context, req domain.UserRegisterRequest) (domain.UserRegisterResponse, error) {
	isRegister := s.userRepository.CheckUserByEmail(ctx, req.Email)
	if isRegister {
		return domain.UserRegisterResponse{}, domain.ErrEmailAlreadyExists
	}
	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return domain.UserRegisterResponse{}, err
	}

	var profilePicture, headline string
	userID := uuid.New()
	allowedMimetype := []string{"image/jpeg", "image/png"}
	if req.ProfilePicture != nil {
		filename := fmt.Sprintf("ProfilePicture-%s", userID)
		objectKey, err := s.awsS3.UploadFile(filename, req.ProfilePicture, "profile-picture", allowedMimetype...)
		if err != nil {
			return domain.UserRegisterResponse{}, domain.ErrUploadFile
		}
		profilePicture = s.awsS3.GetPublicLinkKey(objectKey)
	}

	if req.Headline != nil {
		filename := fmt.Sprintf("Headline-%s", userID)
		objectKey, err := s.awsS3.UploadFile(filename, req.Headline, "headline", allowedMimetype...)
		if err != nil {
			return domain.UserRegisterResponse{}, domain.ErrUploadFile
		}
		headline = s.awsS3.GetPublicLinkKey(objectKey)
	}

	user := entities.User{
		ID:             userID,
		Name:           req.Name,
		Password:       password,
		Email:          req.Email,
		About:          req.About,
		Address:        req.Address,
		CurrentTitle:   req.CurrentTitle,
		ProfilePicture: profilePicture,
		Headline:       headline,
		IsPremium:      false,
		Role:           domain.RoleUser,
		Slug:           utils.CreateSlug(req.Name),
	}

	create, err := s.userRepository.RegisterUser(ctx, user)
	if err != nil {
		return domain.UserRegisterResponse{}, domain.ErrRegisterUserFailed
	}
	return domain.UserRegisterResponse{
		Name:           create.Name,
		Email:          create.Email,
		About:          create.About,
		Address:        create.Address,
		CurrentTitle:   create.CurrentTitle,
		ProfilePicture: create.ProfilePicture,
		Headline:       create.Headline,
		IsPremium:      create.IsPremium,
	}, nil
}

func (s *userService) Login(ctx context.Context, req domain.UserLoginRequest) (domain.UserLoginResponse, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return domain.UserLoginResponse{}, domain.ErrUserNotFound
	}
	if ok := utils.CheckPassword(req.Password, user.Password); !ok {
		return domain.UserLoginResponse{}, domain.CredentialInvalid
	}

	if user.Role != "user" {
		return domain.UserLoginResponse{}, domain.ErrUserNotFound
	}

	token := s.jwtService.GenerateTokenUser(user.ID.String(), user.Role)

	return domain.UserLoginResponse{
		Email:          user.Email,
		Token:          token,
		Role:           user.Role,
		Slug:           user.Slug,
		Name:           user.Name,
		CurrentTitle:   user.CurrentTitle,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, slug string) (domain.UserProfileResponse, error) {
	res, err := s.userRepository.GetProfile(ctx, slug)

	if err != nil {
		return domain.UserProfileResponse{}, domain.ErrGetProfile
	}

	return res, nil

}

func (s *userService) UpdateProfile(ctx context.Context, req domain.UpdateUserRequest, userID string) error {
	user := entities.User{
		Name:         req.Name,
		About:        req.About,
		Address:      req.Address,
		CurrentTitle: req.CurrentTitle,
	}

	allowedMimetype := []string{"image/jpeg", "image/png"}
	if req.ProfilePicture != nil {

		objectKey, err := s.awsS3.UploadFile(utils.GenerateRandomFileName(user.ProfilePicture), req.ProfilePicture, "profile-picture", allowedMimetype...)
		if err != nil {
			return domain.ErrUploadFile
		}

		user.ProfilePicture = s.awsS3.GetPublicLinkKey(objectKey)
	}

	if req.Headline != nil {
		objectKey, err := s.awsS3.UploadFile(utils.GenerateRandomFileName(user.Headline), req.Headline, "headline", allowedMimetype...)
		if err != nil {
			return domain.ErrUploadFile
		}
		user.Headline = s.awsS3.GetPublicLinkKey(objectKey)

	}

	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	if err := s.userRepository.UpdateProfile(ctx, user, parsedUserID); err != nil {
		return err
	}
	return nil
}

func (s *userService) PostEducation(ctx context.Context, req domain.PostUserEducationRequest, userID string) error {
	if exist := s.userRepository.CheckUserByID(ctx, userID); !exist {
		return domain.ErrUserNotFound
	}
	userid, err := uuid.Parse(userID)
	if err != nil {
		return domain.ErrParseUUID
	}

	education := entities.UserEducation{
		ID:           uuid.New(),
		UserID:       userid,
		SchoolName:   req.SchoolName,
		Degree:       req.Degree,
		FieldOfStudy: req.FieldOfStudy,
		Description:  req.Description,
		StartedAt:    utils.ConvertStringToTime(req.StartDate),
		EndedAt:      time.Time{},
	}

	if err := s.userRepository.PostEducation(ctx, education); err != nil {
		return domain.ErrUpdateEducation
	}
	return nil
}

func (s *userService) UpdateEducation(ctx context.Context, req domain.UpdateUserEducationRequest, userID string) error {
	if exist := s.userRepository.CheckUserByID(ctx, userID); !exist {
		return domain.ErrUserNotFound
	}

	userid, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	educationID, err := uuid.Parse(req.EducationID)

	if err != nil {
		return domain.ErrParseUUID
	}

	userEducation := entities.UserEducation{
		ID:           educationID,
		UserID:       userid,
		SchoolName:   req.SchoolName,
		Degree:       req.Degree,
		FieldOfStudy: req.FieldOfStudy,
		Description:  req.Description,
		StartedAt:    utils.ConvertStringToTime(req.StartDate),
		EndedAt:      utils.ConvertStringToTime(req.EndDate),
	}

	if err := s.userRepository.UpdateEducation(ctx, userEducation); err != nil {
		return domain.ErrUpdateExperience
	}

	return nil
}

func (s *userService) DeleteEducation(ctx context.Context, educationID string) error {
	// if exist := s.userRepository.CheckUserByID(ctx, userID
	// ); !exist {
	// 	return domain.ErrUserNotFound

	id, err := uuid.Parse(educationID)

	if err != nil {
		return domain.ErrParseUUID
	}

	if err := s.userRepository.DeleteEducation(ctx, id); err != nil {
		return domain.ErrDeleteEducation
	}

	return nil
}

func (s *userService) PostExperience(ctx context.Context, req domain.PostUserExperienceRequest, userID string) error {
	if exist := s.userRepository.CheckUserByID(ctx, userID); !exist {
		return domain.ErrUserNotFound
	}

	userid, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	companyID, err := uuid.Parse(req.CompanyID)

	if err != nil {
		return domain.ErrParseUUID
	}

	userExperience := entities.UserExperience{
		ID:          uuid.New(),
		UserID:      userid,
		Title:       req.Title,
		CompanyID:   companyID,
		Location:    req.Location,
		Description: req.Description,
		StartedAt:   utils.ConvertStringToTime(req.StartDate),
		EndedAt:     utils.ConvertStringToTime(req.EndDate),
	}

	if err := s.userRepository.PostExperience(ctx, userExperience); err != nil {
		return domain.ErrPostExperience
	}

	return nil

}

func (s *userService) UpdateExperience(ctx context.Context, req domain.UpdateUserExperienceRequest, userID string) error {
	// if exist := s.userRepository.CheckUserByID(ctx, userID
	// ); !exist {
	// 	return domain.ErrUserNotFound
	// }

	userid, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	companyID, err := uuid.Parse(req.CompanyID)

	if err != nil {
		return domain.ErrParseUUID
	}

	experienceID, err := uuid.Parse(req.ExperienceID)

	if err != nil {
		return domain.ErrParseUUID
	}

	userExperience := entities.UserExperience{
		ID:          experienceID,
		UserID:      userid,
		Title:       req.Title,
		CompanyID:   companyID,
		Location:    req.Location,
		Description: req.Description,
		StartedAt:   utils.ConvertStringToTime(req.StartDate),
		EndedAt:     utils.ConvertStringToTime(req.EndDate),
	}

	if err := s.userRepository.UpdateExperience(ctx, userExperience); err != nil {
		return domain.ErrUpdateExperience
	}

	return nil
}

func (s *userService) DeleteExperience(ctx context.Context, experienceID string) error {

	id, err := uuid.Parse(experienceID)

	if err != nil {
		return domain.ErrParseUUID
	}

	if err := s.userRepository.DeleteExperience(ctx, id); err != nil {
		return domain.ErrDeleteExperience
	}

	return nil
}

func (s *userService) PostSkill(ctx context.Context, req domain.PostUserSkillRequest, userID string) error {
	if exist := s.userRepository.CheckUserByID(ctx, userID); !exist {
		return domain.ErrUserNotFound
	}

	userid, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	skillid, err := uuid.Parse(req.SkillID)

	if err != nil {
		return domain.ErrParseUUID
	}

	skill := entities.UserSkill{
		ID:      uuid.New(),
		UserID:  userid,
		SkillID: skillid,
	}

	if err := s.userRepository.PostSkill(ctx, skill); err != nil {
		return domain.ErrPostSkill
	}

	return nil
}

func (s *userService) DeleteSkill(ctx context.Context, skillID string) error {
	// if exist := s.userRepository.CheckUserByID(ctx, userID
	// ); !exist {
	// 	return domain.ErrUserNotFound

	id, err := uuid.Parse(skillID)

	if err != nil {
		return domain.ErrParseUUID
	}

	if err := s.userRepository.DeleteSkill(ctx, id); err != nil {
		return domain.ErrDeleteSkill
	}

	return nil
}

func (s *userService) SearchUser(ctx context.Context, query domain.UserSearchRequest) ([]domain.UserSearchResponse, error) {
	var usersResponse []domain.UserSearchResponse

	users, err := s.userRepository.SearchUser(ctx, query)

	if err != nil {
		return nil, domain.ErrSearchUser
	}

	for _, user := range users {
		usersResponse = append(usersResponse, domain.UserSearchResponse{
			ID:             user.ID.String(),
			Name:           user.Name,
			Slug:           user.Slug,
			Type:           user.Role,
			ProfilePicture: user.ProfilePicture,
			Headline:       user.CurrentTitle,
		})
	}

	if usersResponse == nil {
		usersResponse = []domain.UserSearchResponse{}
	}

	return usersResponse, nil
}

func (s *userService) GetSkills(ctx context.Context) ([]domain.SkillsResponse, error) {
	var skillsResponse []domain.SkillsResponse

	skills, err := s.userRepository.GetSkills(ctx)

	if err != nil {
		return nil, domain.ErrGetSkills
	}

	for _, skill := range skills {
		skillsResponse = append(skillsResponse, domain.SkillsResponse{
			ID:   skill.ID.String(),
			Name: skill.Name,
		})
	}

	if skillsResponse == nil {
		skillsResponse = []domain.SkillsResponse{}
	}

	return skillsResponse, nil

}
