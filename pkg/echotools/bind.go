package echotools

import "github.com/labstack/echo/v4"

func Bind[T any](ctx echo.Context) (*T, error) {
	item := new(T)

	if err := ctx.Bind(item); err != nil {
		return nil, err
	}

	return item, nil
}
