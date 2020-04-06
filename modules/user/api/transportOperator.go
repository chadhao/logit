package api

import (
	"errors"
	"net/http"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/request"
	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TransportOperatorRegister(c echo.Context) error {
	tr := request.TransportOperatorRegRequest{}

	if err := c.Bind(&tr); err != nil {
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	user := &model.User{ID: uid}
	if err := user.Find(); err != nil {
		return errors.New("cannot find user")
	}

	// Assign transport operator super identity
	if _, err := tr.Reg(uid); err != nil {
		return err
	}

	if !utils.RolesAssert(user.RoleIDs).Is(constant.ROLE_TO_SUPER) {
		// Update user role
		user.RoleIDs = append(user.RoleIDs, constant.ROLE_TO_SUPER)
		if err := user.Update(); err != nil {
			return err
		}
	}

	// Issue token
	token, err := user.IssueToken(c.Get("config").(config.Config))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

func TransportOperatorApply(c echo.Context) error {
	r := struct {
		TransportOperatorID string `json:"transportOperatorID" query:"transportOperatorID"`
	}{}

	if err := c.Bind(&r); err != nil {
		return err
	}
	toID, err := primitive.ObjectIDFromHex(r.TransportOperatorID)
	if err != nil {
		return err
	}
	to := &model.TransportOperator{
		ID: toID,
	}
	if err := to.Find(); err != nil {
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	identity, err := to.AddIdentity(uid, model.TO_DRIVER, "")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, identity)
}

func TransportOperatorUpdate(c echo.Context) error {
	tr := request.TransportOperatorUpdateRequest{}

	if err := c.Bind(&tr); err != nil {
		return err
	}
	uid, _ := c.Get("user").(primitive.ObjectID)
	ti := model.TransportOperatorIdentity{
		UserID:              uid,
		TransportOperatorID: tr.ID,
		Identity:            model.TO_SUPER,
	}
	if tos, err := ti.Filter(); len(tos) < 1 || err != nil {
		return errors.New("no authorization")
	}

	to, err := tr.Update()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, to)
}

func TransportOperatorAddIdentity(c echo.Context) error {
	tr := request.TransportOperatorAddIdentityRequest{}

	if err := c.Bind(&tr); err != nil {
		return err
	}
	uid, _ := c.Get("user").(primitive.ObjectID)
	if !model.IsIdentity(uid, tr.TransportOperatorID, []model.TOIdentity{model.TO_SUPER}) {
		return errors.New("no authorization")
	}

	assignedUser := &model.User{ID: tr.UserID}
	if err := assignedUser.Find(); err != nil {
		return errors.New("cannot find user")
	}

	identity, err := tr.Add()
	if err != nil {
		return err
	}

	// update user role
	roles := utils.RolesAssert(assignedUser.RoleIDs)
	if tr.Identity == model.TO_SUPER && !roles.Is(constant.ROLE_TO_SUPER) {
		assignedUser.RoleIDs = append(assignedUser.RoleIDs, constant.ROLE_TO_SUPER)
		if err := assignedUser.Update(); err != nil {
			return err
		}
	}
	if tr.Identity == model.TO_ADMIN && !roles.Is(constant.ROLE_TO_ADMIN) {
		assignedUser.RoleIDs = append(assignedUser.RoleIDs, constant.ROLE_TO_ADMIN)
		if err := assignedUser.Update(); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, identity)
}

func TransportOperatorRemoveIdentity(c echo.Context) error {
	r := struct {
		TransportOperatorIdentityID string `json:"id" query:"id"`
	}{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	rid, err := primitive.ObjectIDFromHex(r.TransportOperatorIdentityID)
	if err != nil {
		return err
	}
	uid, _ := c.Get("user").(primitive.ObjectID)

	toi := model.TransportOperatorIdentity{
		ID: rid,
	}
	if err := toi.Find(); err != nil {
		return err
	}

	if !model.IsIdentity(uid, toi.TransportOperatorID, []model.TOIdentity{model.TO_SUPER}) {
		return errors.New("no authorization")
	}

	if toi.UserID == uid {
		return errors.New("cannot delete")
	}

	if err := toi.Delete(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "ok")
}
