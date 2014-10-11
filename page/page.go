package page

import (
	"fmt"
	"github.com/sclevine/agouti/page/internal/webdriver"
	"image/png"
	"io"
	"os"
)

type Page struct {
	Driver driver
}

type driver interface {
	Navigate(url string) error
	GetElements(selector string) ([]webdriver.Element, error)
	GetWindow() (webdriver.Window, error)
	Screenshot() (io.Reader, error)
	SetCookie(cookie *webdriver.Cookie) error
	DeleteCookie(name string) error
	DeleteAllCookies() error
	GetURL() (string, error)
}

func (p *Page) Navigate(url string) error {
	if err := p.Driver.Navigate(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (p *Page) SetCookie(name string, value interface{}, path, domain string, secure, httpOnly bool, expiry int64) error {
	cookie := webdriver.Cookie{name, value, path, domain, secure, httpOnly, expiry}
	if err := p.Driver.SetCookie(&cookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
	}
	return nil
}

func (p *Page) DeleteCookie(name string) error {
	if err := p.Driver.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %s", name, err)
	}
	return nil
}

func (p *Page) ClearCookies() error {
	if err := p.Driver.DeleteAllCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %s", err)
	}
	return nil
}

func (p *Page) URL() (string, error) {
	url, err := p.Driver.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %s", err)
	}
	return url, nil
}

func (p *Page) Size(width, height int) error {
	window, err := p.Driver.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to retrieve window: %s", err)
	}

	if err := window.SetSize(width, height); err != nil {
		return fmt.Errorf("failed to set window size: %s", err)
	}

	return nil
}

func (p *Page) Screenshot(filepath, filename string) error {
	if err := os.MkdirAll(filepath, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	if err := os.Chdir(filepath); err != nil {
		return fmt.Errorf("failed to switch directories: %s", err)
	}

	screenshot, err := p.Driver.Screenshot()
	if err != nil {
		return fmt.Errorf("failed to retrieve screenshot: %s", err)
	}

	decodedImage, err := png.Decode(screenshot)
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %s", err)
	}

	file, err := os.Create(filename + ".png")
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	if err = png.Encode(file, decodedImage); err != nil {
		file.Close()
		return fmt.Errorf("failed to save PNG: %s", err)
	}

	file.Close()
	return nil
}

func (p *Page) Find(selector string) Selection {
	return &selection{p.Driver, []string{selector}}
}

func (p *Page) Selector() string {
	return p.body().Selector()
}

func (p *Page) Click() error {
	return p.body().Click()
}

func (p *Page) Check() error {
	return p.body().Check()
}

func (p *Page) Fill(text string) error {
	return p.body().Fill(text)
}

func (p *Page) Text() (string, error) {
	return p.body().Text()
}

func (p *Page) Attribute(attribute string) (string, error) {
	return p.body().Attribute(attribute)
}

func (p *Page) CSS(property string) (string, error) {
	return p.body().CSS(property)
}

func (p *Page) Selected() (bool, error) {
	return p.body().Selected()
}

func (p *Page) Select(text string) error {
	return p.body().Select(text)
}

func (p *Page) body() *selection {
	return &selection{p.Driver, []string{"body"}}
}
