package model

import (
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(150);unique;check:length(trim(name)) >= 3"`

	// many-to-many (role-to-permission)
	Permissions []Permission `gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:RoleID;references:ID;joinReferences:PermissionID"`
	// many-to-many (role-to-user)
	Users []User `gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:RoleID;references:ID;joinReferences:UserID"`
}

func AddConstraintsRolePermissions(db *gorm.DB) error {
	return db.Exec(`
        ALTER TABLE role_permissions 
		    DROP CONSTRAINT IF EXISTS fk_role_permissions_role,
		    DROP CONSTRAINT IF EXISTS fk_role_permissions_permission,
			ADD CONSTRAINT fk_role_permissions_role
				FOREIGN KEY (role_id) REFERENCES roles(id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT,
			ADD CONSTRAINT fk_role_permissions_permission
				FOREIGN KEY (permission_id) REFERENCES permissions(id)
				ON DELETE CASCADE
				ON UPDATE RESTRICT;
    `).Error
}
