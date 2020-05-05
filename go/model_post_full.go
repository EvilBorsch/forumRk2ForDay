/*
 * forum
 *
 * Тестовое задание для реализации проекта \"Форумы\" на курсе по базам данных в Технопарке Mail.ru (https://park.mail.ru).
 *
 * API version: 0.1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import post "go-server-server-generated/src/post/models"

// Полная информация о сообщении, включая связанные объекты.
type PostFull struct {
	Post *post.Post `json:"post,omitempty"`

	Author *User `json:"author,omitempty"`

	Thread *Thread `json:"thread,omitempty"`

	Forum *Forum `json:"forum,omitempty"`
}