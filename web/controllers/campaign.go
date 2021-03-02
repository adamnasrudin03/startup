package controllers

import (
	"bwa-startup/campaign"
	"bwa-startup/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

type campaignController struct {
	campaignService campaign.Service
	userService users.Service
}

func NewcampaignController(campaignService campaign.Service, userService users.Service) *campaignController {
	return &campaignController{campaignService, userService}
}

func (h *campaignController) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignController) New(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := campaign.FormCreateCampaignInput{}
	input.Users = users

	c.HTML(http.StatusOK, "campaign_new.html", input)
}