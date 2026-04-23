package minio

//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/joho/godotenv"
//	"log"
//	"minio-gin-crud/internal/common/config"
//	"minio-gin-crud/pkg/minio"
//)
//
//func run() {
//	// Загрузка конфигурации из файла .env
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatalf("Ошибка загрузки файла .env: %v", err)
//	}
//
//	// Инициализация соединения с Minio
//	minioClient := minio.NewMinioClient()
//	err = minioClient.InitMinio()
//	if err != nil {
//		log.Fatalf("Ошибка инициализации Minio: %v", err)
//	}
//
//	// Инициализация маршрутизатора Gin
//	router := gin.Default()
//
//	// Запуск сервера Gin
//	port := config.AppConfig.Port // Мы берем
//	err = router.Run(":" + port)
//	if err != nil {
//		log.Fatalf("Ошибка запуска сервера Gin: %v", err)
//	}
//}
