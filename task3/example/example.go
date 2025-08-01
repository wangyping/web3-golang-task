package example

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type Student struct {
	ID        uint      // Standard field for the primary key
	Name      string    // A regular string field
	Age       uint8     // An unsigned 8-bit integer
	Grade     string    // A pointer to time.Time, can be null
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time // Automatically managed by GORM for update time
	ignored   string    // fields that aren't exported are ignored
}

type Account struct {
	ID          uint
	AccountName string
	Balance     float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Transactions struct {
	ID            uint
	Amount        float64
	FromAccountId uint
	ToAccountId   uint
}

// Employee 结构体定义，用于映射 employees 表
type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

// Book 结构体定义
type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

type User struct {
	gorm.Model
	Name      string
	Posts     []Post
	PostCount int
}

type Post struct {
	gorm.Model
	Title         string
	Comments      []Comment
	UserID        uint
	CommentStatus string
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

func Run(db *gorm.DB) {
	// 题目1：基本CRUD操作
	// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
	// 要求 ：
	// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	// db.AutoMigrate(&Student{})
	// stu := Student{Name: "张三", Age: 20, Grade: "三年级"}
	// db.Debug().Create(&stu)
	// db.Debug().Where("age > ?", 18).Find(&Student{})
	// db.Debug().Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")
	// db.Debug().Where("age < ?", 15).Delete(&Student{})

	// 题目2：事务语句
	// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
	// 要求 ：
	// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
	// db.Debug().AutoMigrate(&Account{})
	// db.Debug().AutoMigrate(&Transactions{})
	// acc := []Account{{ID: 1, AccountName: "A", Balance: 1000}, {ID: 2, AccountName: "B", Balance: 1000}}
	// db.Debug().Create(&acc)
	// err := db.Transaction(func(tx *gorm.DB) error {
	// 	var accounts []Account
	// 	if err := tx.Debug().Where("account_name in ?", []string{"A", "B"}).Find(&accounts).Error; err != nil {
	// 		return err
	// 	}
	// 	if len(accounts) != 2 {
	// 		return errors.New("no match account")
	// 	}
	// 	m := make(map[string]Account)
	// 	for _, account := range accounts {
	// 		m[account.AccountName] = account
	// 	}
	// 	if m["A"].Balance < 100 {
	// 		return errors.New("Account a balance is not enough")
	// 	}
	// 	if err := tx.Debug().Model(&Account{}).Where("account_name = ?", "A").Update("balance", gorm.Expr("balance - ?", 100)).Error; err != nil {
	// 		return err
	// 	}
	// 	if err := tx.Debug().Model(&Account{}).Where("account_name = ?", "B").Update("balance", gorm.Expr("balance + ?", 100)).Error; err != nil {
	// 		return err
	// 	}

	// 	tx.Debug().Create(&Transactions{FromAccountId: m["A"].ID, ToAccountId: m["B"].ID, Amount: 100})
	// 	// 在这里可以处理查询结果
	// 	return nil
	// })
	// if err != nil {
	// 	fmt.Println("事务处理失败:", err)
	// } else {
	// 	fmt.Println("事务处理成功")
	// }

	// db.AutoMigrate(&Employee{})
	// db.AutoMigrate(&Book{})
	// db.Create(&Employee{Name: "张三", Department: "技术部", Salary: 5000})
	// db.Create(&Employee{Name: "李四", Department: "销售部", Salary: 4000})
	// db.Create(&Employee{Name: "王五", Department: "财务部", Salary: 3000})
	// db.Create(&Employee{Name: "孙六", Department: "技术部", Salary: 5300})

	// db.Create(&Book{Title: "Go 语言基础", Author: "小王子", Price: 28.8})
	// db.Create(&Book{Title: "Go 语言进阶", Author: "小王子", Price: 38.8})
	// db.Create(&Book{Title: "Go 语言实战", Author: "小王子", Price: 48.8})
	// db.Create(&Book{Title: "Go 语言微服务", Author: "小王子", Price: 58.8})

	// 进阶gorm
	// 题目1：模型定义
	// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
	// 要求 ：
	// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
	// 编写Go代码，使用Gorm创建这些模型对应的数据库表。
	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&Post{})
	// db.AutoMigrate(&Comment{})
	// db.Debug().Create(&User{Name: "张三", Posts: []Post{{Title: "第一篇博客", Comments: []Comment{{Content: "好"}, {Content: "cai"}}}, {Title: "第二篇博客", Comments: []Comment{{Content: "11"}, {Content: "222"}}}}})
	// db.Debug().Create(&User{Name: "里斯", Posts: []Post{{Title: "demo1博客", Comments: []Comment{{Content: "dd"}, {Content: "cc"}}}, {Title: "demo2博客", Comments: []Comment{{Content: "33"}, {Content: "44"}}}}})

	// 题目2：关联查询
	// 基于上述博客系统的模型定义。
	// 要求 ：
	// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	// 编写Go代码，使用Gorm查询评论数量最多的文章信息。

	user := User{}
	db.Debug().Preload("Posts").Preload("Posts.Comments").Where("name = ?", "张三").Find(&user)
	fmt.Printf("%+v\n", user)

	posts := []Post{}
	//comments := []Comment{}
	db.Debug().Preload("Comments").Find(&posts)
	for _, post := range posts {
		fmt.Println(post.Title, len(post.Comments))
	}

	var post Post
	result := db.Debug().
		Preload("Comments").
		Joins("LEFT JOIN comments ON posts.id = comments.post_id").
		Group("posts.id").
		Order("COUNT(comments.id) DESC").
		Limit(1).
		First(&post)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("没有找到文章")
		} else {
			fmt.Printf("查询失败: %v\n", result.Error)
		}
		return
	}

	fmt.Printf("评论数量最多的文章: %s\n", post.Title)
	fmt.Printf("评论数量: %d\n", len(post.Comments))

	// 输出所有评论内容
	for _, comment := range post.Comments {
		fmt.Printf("  评论: %s\n", comment.Content)
	}

	// 题目3：钩子函数
	// 继续使用博客系统的模型。
	// 要求 ：
	// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
	comment := Comment{Model: gorm.Model{ID: 1}, PostID: 1}
	db.Debug().Delete(&comment)
}

// Post模型的钩子函数，在创建后自动更新用户的文章数量
func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 检查UserID是否有效
	if p.UserID > 0 {
		// 使用表达式直接更新计数，避免并发问题
		if err := tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// post := Post{}
	// tx.Model(&Post{}).Joins("LEFT JOIN comments ON posts.id = comments.post_id").Group("posts.id").Where("id = ? ", c.PostID).Having("COUNT(id)=1").First(&post)
	// fmt.Println(post)
	// fmt.Println(c.PostID)
	// 检查PostID是否有效
	if c.PostID > 0 {
		// 使用表达式直接更新计数，避免并发问题
		if err := tx.Model(&Post{}).Joins("LEFT JOIN comments ON posts.id = comments.post_id").Group("posts.id").Where("id = ? ", c.PostID).Having("COUNT(id)=1").UpdateColumn("comment_status", "无评论").Error; err != nil {
			return err
		}
	}
	return nil
}

func RunSqlX(db *sqlx.DB) {
	// Sqlx入门
	// 题目1：使用SQL扩展库进行查询
	// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	// 要求 ：
	// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	// 题目2：实现类型安全映射
	// 假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	// 要求 ：
	// 定义一个 Book 结构体，包含与 books 表对应的字段。
	// 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
	// 数据库连接字符串
	// 查询所有部门为"技术部"的员工信息
	var techEmployees []Employee
	err := db.Select(&techEmployees, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")
	if err != nil {
		fmt.Printf("查询技术部员工失败: %v\n", err)
	} else {
		fmt.Printf("技术部员工数量: %d\n", len(techEmployees))
		for _, emp := range techEmployees {
			fmt.Printf("员工: %+v\n", emp)
		}
	}

	var maxSalaryEmployee Employee
	err = db.Get(&maxSalaryEmployee, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		fmt.Printf("查询最高工资员工失败: %v\n", err)
	} else {
		fmt.Printf("工资最高的员工: %+v\n", maxSalaryEmployee)
	}

	// 查询价格大于50元的书籍
	var expensiveBooks []Book
	err = db.Select(&expensiveBooks, "SELECT id, title, author, price FROM books WHERE price > ?", 50)
	if err != nil {
		fmt.Printf("查询高价书籍失败: %v\n", err)
	} else {
		fmt.Printf("价格大于50元的书籍数量: %d\n", len(expensiveBooks))
		for _, book := range expensiveBooks {
			fmt.Printf("书籍: %+v\n", book)
		}
	}
}
