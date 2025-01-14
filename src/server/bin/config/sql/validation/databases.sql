-- Dwarf Wars Database (World tracking, directory tracking, ...)
CREATE DATABASE IF NOT EXISTS DW 
	DEFAULT CHARACTER SET = 'utf8mb4'
	DEFAULT COLLATE = 'utf8mb4_unicode_ci'
;

-- Dwarf Wars Security Database (Accounts, SSO, 2FA, ...)
CREATE DATABASE IF NOT EXISTS DWS
	DEFAULT CHARACTER SET = 'utf8mb4'
	DEFAULT COLLATE = 'utf8mb4_unicode_ci'
;

