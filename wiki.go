package main

import "os"

type Page struct {
	Title string //标题
	Body  []byte //正文 类型定位为[]byte类型 而不是string，为了方便存储
}

//保存Page到一个文本文件
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600) //0600让当前用户拥有读写权限
}

//将保存的Page读取出来
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
