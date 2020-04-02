/*
 * forum
 *
 * Тестовое задание для реализации проекта \"Форумы\" на курсе по базам данных в Технопарке Mail.ru (https://park.mail.ru). 
 *
 * API version: 0.1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type status struct {

	// Кол-во пользователей в базе данных.
	User float32 `json:"user"`

	// Кол-во разделов в базе данных.
	Forum float32 `json:"forum"`

	// Кол-во веток обсуждения в базе данных.
	Thread float32 `json:"thread"`

	// Кол-во сообщений в базе данных.
	Post float32 `json:"post"`
}
