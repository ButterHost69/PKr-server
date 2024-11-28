release:
	go run . -r

debug:
	go run .

testserver:
	go run . -t

# FIXME: [ ] Not working all the time
deletetest:
	if (Test-Path "test_database.db") { Remove-Item "test_database.db"; Write-Host "File 'test_database.db' deleted." } else { Write-Host "File 'test_database.db' does not exist." }