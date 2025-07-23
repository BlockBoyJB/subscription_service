package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"subscription_service/internal/service"
	"time"
)

type subscriptionRouter struct {
	sub service.Subscription
}

func newSubscriptionRouter(g *echo.Group, sub service.Subscription) {
	r := &subscriptionRouter{
		sub: sub,
	}

	g.POST("", r.create)
	g.GET("/all", r.findAll)
	g.GET("/:id", r.findById)
	g.GET("/price", r.findPrice)
	g.PUT("/:id", r.update)
	g.DELETE("/:id", r.delete)
}

type subscriptionInput struct {
	ServiceName string  `json:"service_name" validate:"required"`
	Price       int     `json:"price" validate:"required"`
	UserId      string  `json:"user_id" validate:"required,uuid4"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     *string `json:"end_date"`
}

// @Summary		Create
// @Description	Create new subscription in database
// @Tags			subscription
// @Accept			json
// @Produce		json
// @Param			input	body		subscriptionInput	true	"input"
// @Success		200		{string}	string				"OK"
// @Failure		400		{string}	string				"Bad Request"
// @Failure		500		{string}	string				"Internal Server Error"
// @Router			/api/v1/subscription [post]
func (r *subscriptionRouter) create(c echo.Context) error {
	var input subscriptionInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := c.Validate(&input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	s, err := parseInputDate(input)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = r.sub.Create(c.Request().Context(), s); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// @Summary		Find All
// @Description	Find all subscription in database
// @Tags			subscription
// @Accept			json
// @Produce		json
// @Success		200	{array}		service.SubscriptionInput
// @Failure		400	{string}	string	"Bad Request"
// @Failure		500	{string}	string	"Internal Server Error"
// @Router			/api/v1/subscription/all [get]
func (r *subscriptionRouter) findAll(c echo.Context) error {
	s, err := r.sub.FindAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, s)
}

// @Summary		Find by id
// @Description	Find subscription in database by id
// @Tags			subscription
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"id"
// @Success		200	{object}	service.SubscriptionOutput
// @Failure		400	{string}	string	"Bad Request"
// @Failure		404	{string}	string	"Not Found"
// @Failure		500	{string}	string	"Internal Server Error"
// @Router			/api/v1/subscription/{id} [get]
func (r *subscriptionRouter) findById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	s, err := r.sub.FindById(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, s)
}

type subscriptionPriceOutput struct {
	Price int `json:"price"`
}

// @Summary		Price
// @Description	Find total price for subscriptions for time interval
// @Tags			subscription
// @Accept			json
// @Produce		json
// @Param			service_name	query		string	false	"name of subscription service"
// @Param			user_id			query		string	false	"user id"
// @Param			start			query		string	true	"start of the time interval. Must be in format mm-yyyy"
// @Param			end				query		string	true	"end of the time interval. Must be in format mm-yyyy"
// @Success		200				{object}	subscriptionPriceOutput
// @Failure		400				{string}	string	"Bad Request"
// @Failure		500				{string}	string	"Internal Server Error"
// @Router			/api/v1/subscription/price [get]
func (r *subscriptionRouter) findPrice(c echo.Context) error {
	start, err := time.Parse("01-2006", c.QueryParam("start"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	end, err := time.Parse("01-2006", c.QueryParam("end"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	price, err := r.sub.FindPrice(c.Request().Context(), service.PriceInput{
		ServiceName: c.QueryParam("service_name"),
		UserId:      c.QueryParam("user_id"),
		StartDate:   start,
		EndDate:     end,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, subscriptionPriceOutput{
		Price: price,
	})
}

// @Summary		Update
// @Description	Update subscription in database by id
// @Tags			subscription
// @Accept			json
// @Produce		json
// @Param			id		path		int					true	"id"
// @Param			input	body		subscriptionInput	true	"input"
// @Success		200		{string}	string				"OK"
// @Failure		400		{string}	string				"Bad Request"
// @Failure		404		{string}	string				"Not Found"
// @Failure		500		{string}	string				"Internal Server Error"
// @Router			/api/v1/subscription/{id} [put]
func (r *subscriptionRouter) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var input subscriptionInput

	if err = c.Bind(&input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err = c.Validate(&input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	s, err := parseInputDate(input)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = r.sub.Update(c.Request().Context(), id, s); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// @Summary		Delete
// @Description	Delete subscription in database by id
// @Tags			subscription
// @Accept			json
// @Produce		json
// @Param			id	path		int		true	"id"
// @Success		200	{string}	string	"OK"
// @Failure		400	{string}	string	"Bad Request"
// @Failure		404	{string}	string	"Not Found"
// @Failure		500	{string}	string	"Internal Server Error"
// @Router			/api/v1/subscription/{id} [delete]
func (r *subscriptionRouter) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = r.sub.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func parseInputDate(input subscriptionInput) (service.SubscriptionInput, error) {
	start, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return service.SubscriptionInput{}, err
	}
	s := service.SubscriptionInput{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserId:      input.UserId,
		StartDate:   start,
	}
	if input.EndDate != nil {
		end, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			return service.SubscriptionInput{}, err
		}
		s.EndDate = &end
	}
	return s, nil
}
