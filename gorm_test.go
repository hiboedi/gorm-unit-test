package golang_gorm

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:@tcp(localhost:3306)/golang_gorm2?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

// execute sql
func TestExecuteSQL(t *testing.T) {
	err := db.Exec("insert into sample(id,name) values (?,?)", "1", "Budi").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id,name) values (?,?)", "2", "Fatir").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id,name) values (?,?)", "3", "Abdur").Error
	assert.Nil(t, err)

	err = db.Exec("insert into sample(id,name) values (?,?)", "4", "Eko").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   string
	Name string
}

func TestRawSQL(t *testing.T) {
	var sample Sample
	err := db.Raw("select id, name from sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Budi", sample.Name)

	var samples []Sample
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(samples))

}

func TestSQLRow(t *testing.T) {
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, Sample{
			Id:   id,
			Name: name,
		})
	}
	assert.Equal(t, 4, len(samples))
}

func TestScanRow(t *testing.T) {
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}
	assert.Equal(t, 4, len(samples))
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID:       "1",
		Password: "knok",
		Name: Name{
			FirstName:  "Budi",
			MiddleName: "Abdurahman",
			LastName:   "Fatir",
		},
		Information: "Ini di ignore",
	}

	result := db.Create(&user)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(1), result.RowsAffected)
}

func TestBatchInsert(t *testing.T) {
	var users []User
	for i := 2; i < 10; i++ {
		users = append(users, User{
			ID:       strconv.Itoa(i),
			Password: "rahasia",
			Name: Name{
				FirstName: "User " + strconv.Itoa(i),
			},
		})
	}

	result := db.Create(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(8), result.RowsAffected)
}

func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       "10",
			Password: "rahasia",
			Name:     Name{FirstName: "User 10"},
		}).Error

		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       "11",
			Password: "rahasia",
			Name:     Name{FirstName: "User 11"},
		}).Error

		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       "12",
			Password: "rahasia",
			Name:     Name{FirstName: "User 12"},
		}).Error

		if err != nil {
			return err
		}

		return nil
	})

	assert.Nil(t, err)
}

func TestTransactionError(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       "13",
			Password: "rahasia",
			Name:     Name{FirstName: "User 13"},
		}).Error

		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       "11",
			Password: "rahasia",
			Name:     Name{FirstName: "User 11"},
		}).Error

		if err != nil {
			return err
		}

		return nil
	})

	assert.NotNil(t, err)
}

func TestManuakTransactionSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "13",
		Password: "rahasia",
		Name:     Name{FirstName: "User 13"},
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "14",
		Password: "rahasia",
		Name:     Name{FirstName: "User 14"},
	}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestManuakTransactionError(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{
		ID:       "15",
		Password: "rahasia",
		Name:     Name{FirstName: "User 15"},
	}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{
		ID:       "14",
		Password: "rahasia",
		Name:     Name{FirstName: "User 14"},
	}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestQuerySingleObject(t *testing.T) {
	user := User{}
	err := db.First(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)

	user = User{}
	user = User{}
	err = db.Last(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, "9", user.ID)
}

func TestQueryInlineCondition(t *testing.T) {
	user := User{}
	err := db.Take(&user, "id = ?", "5").Error
	assert.Nil(t, err)
	assert.Equal(t, "5", user.ID)
	assert.Equal(t, "User 5", user.Name.FirstName)
}

func TestQueryAllObject(t *testing.T) {
	var users []User
	err := db.Find(&users, "id in ?", []string{"1", "2", "3", "6", "8"}).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User
	err := db.Where("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len((users)))
}

func TestOrOperator(t *testing.T) {
	var users []User
	err := db.Where("first_name = ?", "Budi").Or("password like ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len((users)))
}

func TestNotOperator(t *testing.T) {
	var users []User
	err := db.Not("first_name like ?", "%User%").Not("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len((users)))
}

func TestSelectField(t *testing.T) {
	var users []User
	err := db.Select("id", "first_name").Find(&users).Error
	assert.Nil(t, err)

	for _, user := range users {
		assert.NotNil(t, user.ID)
		assert.NotEqual(t, "", user.Name.FirstName)
	}
	assert.Equal(t, 14, len(users))
}

func TestStructCondition(t *testing.T) {
	userCondition := User{
		Name: Name{
			FirstName: "User 8",
		},
		Password: "rahasia",
	}
	var users []User
	err := db.Where(userCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestMapCondition(t *testing.T) {
	mapCondition := map[string]interface{}{
		"middle_name": "",
		"last_name":   "",
	}
	var users []User
	err := db.Where(mapCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 13, len(users))
}

func TestOrderLimitOffset(t *testing.T) {
	var users []User
	err := db.Order("id asc, first_name desc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse
	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 14, len(users))
	fmt.Println(users)
}

func TestUpdate(t *testing.T) {
	user := User{}
	err := db.Take(&user, "id = ?", "1").Error

	assert.Nil(t, err)

	user.Name.FirstName = "Renal"
	user.Name.MiddleName = ""
	user.Name.LastName = "Aditya"
	user.Password = "rahasia"

	err = db.Save(&user).Error
	assert.Nil(t, err)

}

func TestUpdateSelectedColumn(t *testing.T) {
	result := db.Model(&User{}).Where("id =?", "1").Updates(map[string]interface{}{
		"middle_name": "",
		"last_name":   "Dika",
	}).Error

	assert.Nil(t, result)

	err := db.Model(&User{}).Where("id = ?", "1").Update("password", "knok").Error

	assert.Nil(t, err)

	err = db.Where("id = ?", "1").Updates(User{
		Name: Name{
			FirstName: "Budi",
			LastName:  "Fatir",
		},
	}).Error
	assert.Nil(t, err)

}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := UserLog{
			UserId: strconv.Itoa(i + 1),
			Action: "Test Action",
		}
		err := db.Create(&userLog).Error
		assert.Nil(t, err)

		assert.NotEqual(t, 0, userLog.ID)
		fmt.Println(userLog.ID)
	}
}

func TestSaveOrUpdate(t *testing.T) {
	userLog := UserLog{
		UserId: "1",
		Action: "Test Action",
	}

	err := db.Save(&userLog).Error //insert
	assert.Nil(t, err)

	userLog.UserId = "2"
	err = db.Save(&userLog).Error //update
	assert.Nil(t, err)

}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	user := User{
		ID: "99",
		Name: Name{
			FirstName: "Trent",
		},
	}

	err := db.Save(&user).Error //insert
	assert.Nil(t, err)

	user.Name.FirstName = "User 99 Updated"
	err = db.Save(&user).Error //update
	assert.Nil(t, err)
}

func TestCoflict(t *testing.T) {
	user := User{
		ID: "88",
		Name: Name{
			FirstName: "Trent",
		},
	}

	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error //insert
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	var user User
	err := db.Take(&user, "id=?", "99").Error
	assert.Nil(t, err)

	err = db.Delete(&user).Error
	assert.Nil(t, err)

	err = db.Delete(&User{}, "id=?", "88").Error
	assert.Nil(t, err)

	err = db.Where("id=?", "77").Delete(&User{}).Error
	assert.Nil(t, err)
}

func TestSoftDelete(t *testing.T) {
	todo := Todo{
		UserId:      "1",
		Title:       "test",
		Description: "desc test",
	}
	err := db.Create(&todo).Error
	assert.Nil(t, err)

	err = db.Delete(&todo).Error
	assert.Nil(t, err)
	assert.NotNil(t, todo.DeletedAt)

	var todos []Todo
	err = db.Find(&todos).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestUncoped(t *testing.T) {
	var todo Todo

	err := db.Unscoped().First(&todo, "id=?", "2").Error
	assert.Nil(t, err)
	fmt.Println(todo)

	err = db.Unscoped().Delete(&todo).Error
	assert.Nil(t, err)

	var todos []Todo
	err = db.Unscoped().Find(&todos).Error
	assert.Nil(t, err)

}

func TestLock(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, "id=?", "1").Error
		if err != nil {
			return err
		}

		user.Name.FirstName = "Jaka"
		user.Name.LastName = "Pitung"
		err = tx.Save(&user).Error

		return err
	})
	assert.Nil(t, err)
}

func TestCreateWallet(t *testing.T) {
	wallet := Wallet{
		ID:      "1",
		UserId:  "1",
		Balance: 1000000,
	}

	err := db.Create(&wallet).Error
	assert.Nil(t, err)
}

//! relasi itu bersifat lazy load maka kita harus mencoba preload /joins

func TestRetrieveRelation(t *testing.T) {
	var user User
	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "1", user.Wallet.ID)
}

func TestJoins(t *testing.T) {
	var user User
	err := db.Model(&User{}).Joins("Wallet").Take(&user, "users.id = ?", "1").Error
	assert.Nil(t, err)

	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "1", user.Wallet.ID)
}

func TestAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       "20",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 20",
		},
		Wallet: Wallet{
			ID:      "20",
			Balance: 1000000,
			UserId:  "20",
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)

}

func TestSkipCreateUpdate(t *testing.T) {
	user := User{
		ID:       "21",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 21",
		},
		Wallet: Wallet{
			ID:      "21",
			Balance: 1000000,
			UserId:  "21",
		},
	}

	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)

}

func TestUserAndAddresses(t *testing.T) {
	user := User{
		ID:       "2",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 2",
		},
		Wallet: Wallet{
			ID:      "2",
			UserId:  "2",
			Balance: 5000000,
		},
		Addresses: []Address{
			{
				UserId:  "2",
				Address: "Bandung",
			},
			{
				UserId:  "2",
				Address: "Padalarang",
			},
		},
	}
	err := db.Save(&user).Error
	assert.Nil(t, err)

}

func TestPreloadJoinOneToMany(t *testing.T) {
	var users []User

	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)

}

func TestTake(t *testing.T) {
	var user User

	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Take(&user, "users.id=?", "50").Error
	assert.Nil(t, err)

}

// ! belongs to many to one
func TestBelongsToAddress(t *testing.T) {
	fmt.Println("Preload")
	var addresses []Address
	err := db.Model(&Address{}).Preload("User").Find(&addresses).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(addresses))

	fmt.Println("Preload")
	addresses = []Address{}
	err = db.Model(&Address{}).Joins("User").Find(&addresses).Error
	assert.Nil(t, err)
}

// ! belongs to one to one
func TestBelongsToWallet(t *testing.T) {
	fmt.Println("Preload")
	var wallets []Wallet
	err := db.Model(&Wallet{}).Preload("User").Find(&wallets).Error
	assert.Nil(t, err)

	fmt.Println("Preload")
	wallets = []Wallet{}
	err = db.Model(&Wallet{}).Joins("User").Find(&wallets).Error
	assert.Nil(t, err)
}

func TestCreateManyToMany(t *testing.T) {
	product := Product{
		ID:    "P001",
		Name:  "Apple",
		Price: 50000,
	}
	err := db.Create(&product).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create((map[string]interface{}{
		"user_id":    "1",
		"product_id": "P001",
	})).Error
	assert.Nil(t, err)

}

func TestPreloadManyToMany(t *testing.T) {
	var product Product
	err := db.Preload("LikeByUsers").Find(&product).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(product.LikeByUsers))

}

func TestPreloadManyToManyUser(t *testing.T) {
	var user User
	err := db.Preload("LikeProducts").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(user.LikeProducts))

}

func TestAssociatinFInd(t *testing.T) {
	// !ambil product
	var product Product

	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	// !cek like by user
	var users []User
	err = db.Model(&product).Where("users.first_name like ?", "User%").Association("LikeByUsers").Find(&users)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestAssosiationAppend(t *testing.T) {
	var user User

	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product

	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikeByUsers").Append(&user)
	assert.Nil(t, err)

}

func TestAssociationReplace(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Take(&user, "id = ?", "1").Error
		assert.Nil(t, err)

		wallet := Wallet{
			ID:      "01",
			UserId:  user.ID,
			Balance: 1000000,
		}
		err = tx.Model(&user).Association("Wallet").Replace(&wallet)

		return err
	}).Error()
	assert.Nil(t, err)
}

func TestAssosiationDelete(t *testing.T) {
	var user User

	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product

	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikeByUsers").Delete(&user)
	assert.Nil(t, err)

}

// ? clear relation
func TestAssosiationClear(t *testing.T) {
	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikeByUsers").Clear()
	assert.Nil(t, err)

}

func TestPreloadingWithCondition(t *testing.T) {
	var user User
	err := db.Preload("Wallet", "balance > ?", 100000).Take(&user, "id=?", "1").Error
	assert.Nil(t, err)
	fmt.Println(user)
}

func TestNestedPreloading(t *testing.T) {
	var wallet Wallet

	err := db.Preload("User.Addresses").Take(&wallet, "id=?", "2").Error
	assert.Nil(t, err)

	fmt.Println(t, wallet)
	fmt.Println(t, wallet.User)
	fmt.Println(t, wallet.User.Addresses)
}

func TestPreloadingAll(t *testing.T) {
	var user User
	err := db.Preload(clause.Associations).Take(&user, "id =?", "1").Error
	assert.Nil(t, err)
}

func TestJoinQuery(t *testing.T) {
	var users []User
	err := db.Joins("join wallets on wallets.user_id=users.id").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))

	users = []User{}
	err = db.Joins("Wallet").Find(&users).Error //! left join
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
}

func TestJoinsWithCondition(t *testing.T) {
	var users []User
	err := db.Joins("join wallets on wallets.user_id=users.id AND wallets.balance > ?", 50000).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))

	users = []User{}
	err = db.Joins("Wallet").Where("Wallet.balance > ?", 50000).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestCount(t *testing.T) {
	var count int64
	err := db.Model(&User{}).Joins("Wallet").Where("Wallet.balance > ?", 50000).Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(4), count)
}

type AggregationResult struct {
	TotalBalance int64
	MinBalance   int64
	MaxBalance   int64
	AvgBalance   float64
}

func TestAggregation(t *testing.T) {
	var result AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance", "min(balance) as min_balance", "max(balance) as max_balance", "avg(balance) as avg_balance").Take(&result).Error

	assert.Nil(t, err)
	assert.Equal(t, int64(12000000), result.TotalBalance)
	assert.Equal(t, int64(1000000), result.MinBalance)
	assert.Equal(t, int64(5000000), result.MaxBalance)
	assert.Equal(t, float64(3000000), result.AvgBalance)
}

func TestGroupByAndHaving(t *testing.T) {
	var result []AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance", "min(balance) as min_balance", "max(balance) as max_balance", "avg(balance) as avg_balance").Joins("User").Group("User.id").Having("sum(balance) > ?", 1000000).Find(&result).Error

	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	var users []User
	err := db.WithContext(ctx).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
}

func BrokeWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance =?", 0)
}

func SultanWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance >?", 2000000)
}

func TestScope(t *testing.T) {
	var wallets []Wallet
	err := db.Scopes(BrokeWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)

	wallets = []Wallet{}
	err = db.Scopes(SultanWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)

}

func TestMigrator(t *testing.T) {
	err := db.Migrator().AutoMigrate(&GuestBook{})
	assert.Nil(t, err)
}

func TestHooks(t *testing.T) {
	user := User{
		Password: "Rahasia",
		Name: Name{
			FirstName: "User 100",
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
	assert.NotEqual(t, "", user.ID)

	fmt.Println(user.ID)
}
