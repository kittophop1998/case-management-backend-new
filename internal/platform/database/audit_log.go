package database

import (
	"case-management/internal/domain/model"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AsyncAuditLogger struct {
	db     *gorm.DB
	ch     chan model.AuditLogs
	wg     sync.WaitGroup
	quitCh chan struct{}
}

func NewAsyncAuditLogger(db *gorm.DB, bufferSize int) *AsyncAuditLogger {
	logger := &AsyncAuditLogger{
		db:     db,
		ch:     make(chan model.AuditLogs, bufferSize),
		quitCh: make(chan struct{}),
	}

	logger.startWorker()
	return logger
}

func (l *AsyncAuditLogger) LogAction(ctx *gin.Context, entry model.AuditLogs) {
	select {
	case l.ch <- entry:
		// OK
	default:
		log.Println("[WARN] audit log queue full, dropping entry")
	}
}

func (l *AsyncAuditLogger) startWorker() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			select {
			case entry := <-l.ch:
				if err := l.saveToDB(entry); err != nil {
					log.Println("[ERROR] Failed to write audit log:", err)
				}
			case <-l.quitCh:
				log.Println("[DEBUG] Quit signal received, stopping worker")
				return
			}
		}
	}()
}

func (l *AsyncAuditLogger) saveToDB(entry model.AuditLogs) error {
	fmt.Println("Saving audit log to DB:", entry)
	jsonMeta, err := json.Marshal(entry.Metadata)
	if err != nil {
		return err
	}

	logModel := model.AuditLogs{
		ID:        entry.ID,
		Action:    entry.Action,
		UserID:    entry.UserID,
		Metadata:  datatypes.JSON(jsonMeta),
		CreatedAt: entry.CreatedAt,
	}

	return l.db.Create(&logModel).Error
}

func (l *AsyncAuditLogger) Shutdown() {
	close(l.quitCh)
	l.wg.Wait()
}
