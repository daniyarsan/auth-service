package server


cfg, err := config.Load()
if err != nil {
log.Fatalf("load config: %v", err)
}