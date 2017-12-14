package state

func testVisitorName() string {
	return "fargle"
}

func testVisitorIp() string {
	return "127.0.0.1"
}

func testVisitorObj() Visitor {
	return Visitor{
		Name: testVisitorName(),
		IP:   testVisitorIp(),
	}
}
