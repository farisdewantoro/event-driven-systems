package logger

func (l *logger) setDefaultFields() {
	l.mu.RLock()
	for k, v := range l.opt.DefaultFields {
		l.logEntry = l.logEntry.WithField(k, v)
	}
	l.mu.RUnlock()
}
