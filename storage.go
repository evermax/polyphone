package main

func (s *server) loginClient(client string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.clients[client] = true
}

func (s *server) logoutClient(client string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.clients[client] = false
}

func (s *server) isClientLoggedIn(client string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.clients[client]
}
