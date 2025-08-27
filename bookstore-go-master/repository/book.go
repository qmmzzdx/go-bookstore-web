package repository

import (
	"bookstore/global"
	"bookstore/model"

	"gorm.io/gorm"
)

// BookDAO 书籍数据访问对象
// 封装了所有与书籍相关的数据库操作
type BookDAO struct {
	db *gorm.DB // GORM数据库连接实例
}

// NewBookDAO 创建新的书籍DAO实例
// 返回:
//
//	*BookDAO - 初始化后的书籍数据访问对象
func NewBookDAO() *BookDAO {
	return &BookDAO{
		db: global.GetDB(), // 从全局变量获取数据库连接
	}
}

// GetAllBooks 获取所有书籍
// 返回:
//
//	[]*model.Book - 书籍对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetAllBooks() ([]*model.Book, error) {
	var books []*model.Book
	// 对应SQL: SELECT * FROM books;
	err := b.db.Find(&books).Error
	return books, err
}

// GetBooksByPage 分页获取书籍（只返回上架状态）
// 参数:
//
//	page - 页码，从1开始
//	pageSize - 每页记录数
//
// 返回:
//
//	[]*model.Book - 当前页的书籍对象切片
//	int64 - 符合条件的总记录数
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetBooksByPage(page, pageSize int) ([]*model.Book, int64, error) {
	var books []*model.Book
	var total int64

	// 获取总数（只统计上架状态的书籍）
	// 对应SQL: SELECT COUNT(*) FROM books WHERE status = 1;
	err := b.db.Model(&model.Book{}).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询（只查询上架状态的书籍）
	// 对应SQL: SELECT * FROM books WHERE status = 1 LIMIT pageSize OFFSET offset;
	offset := (page - 1) * pageSize
	err = b.db.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&books).Error
	return books, total, err
}

// GetBooksByPageForAdmin 分页获取书籍（管理员用，不过滤状态）
// 参数:
//
//	page - 页码，从1开始
//	pageSize - 每页记录数
//
// 返回:
//
//	[]*model.Book - 当前页的书籍对象切片
//	int64 - 总记录数
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetBooksByPageForAdmin(page, pageSize int) ([]*model.Book, int64, error) {
	var books []*model.Book
	var total int64

	// 获取总数（不筛选状态）
	// 对应SQL: SELECT COUNT(*) FROM books;
	err := b.db.Model(&model.Book{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询（不筛选状态）
	// 对应SQL: SELECT * FROM books LIMIT pageSize OFFSET offset;
	offset := (page - 1) * pageSize
	err = b.db.Offset(offset).Limit(pageSize).Find(&books).Error
	return books, total, err
}

// GetBooksByPageForAdminWithSearch 分页获取书籍（管理员用，支持搜索）
// 参数:
//
//	page - 页码，从1开始
//	pageSize - 每页记录数
//	title - 书籍标题搜索关键词
//	author - 作者搜索关键词
//	bookType - 书籍类型筛选条件
//	status - 状态筛选条件（指针类型，nil表示不筛选）
//
// 返回:
//
//	[]*model.Book - 当前页的书籍对象切片
//	int64 - 符合条件的总记录数
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetBooksByPageForAdminWithSearch(page, pageSize int, title, author, bookType string, status *int) ([]*model.Book, int64, error) {
	var books []*model.Book
	var total int64

	// 构建查询条件
	query := b.db.Model(&model.Book{})

	// 添加搜索条件
	if title != "" {
		// 对应SQL条件: title LIKE '%title%'
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if author != "" {
		// 对应SQL条件: author LIKE '%author%'
		query = query.Where("author LIKE ?", "%"+author+"%")
	}
	if bookType != "" {
		// 对应SQL条件: type = bookType
		query = query.Where("type = ?", bookType)
	}
	if status != nil {
		// 对应SQL条件: status = status
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	// 对应SQL: SELECT COUNT(*) FROM books [WHERE conditions...]
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	// 对应SQL: SELECT * FROM books [WHERE conditions...] LIMIT pageSize OFFSET offset;
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&books).Error
	return books, total, err
}

// GetHotBooks 获取热销书籍（只返回上架状态）
// 参数:
//
//	limit - 返回的记录数量限制
//
// 返回:
//
//	[]*model.Book - 热销书籍对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetHotBooks(limit int) ([]*model.Book, error) {
	var books []*model.Book
	// 对应SQL: SELECT * FROM books WHERE status = 1 ORDER BY sale DESC LIMIT limit;
	err := b.db.Where("status = ?", 1).Order("sale DESC").Limit(limit).Find(&books).Error
	return books, err
}

// GetNewBooks 获取新书（只返回上架状态）
// 参数:
//
//	limit - 返回的记录数量限制
//
// 返回:
//
//	[]*model.Book - 新书对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetNewBooks(limit int) ([]*model.Book, error) {
	var books []*model.Book
	// 对应SQL: SELECT * FROM books WHERE status = 1 ORDER BY created_at DESC LIMIT limit;
	err := b.db.Where("status = ?", 1).Order("created_at DESC").Limit(limit).Find(&books).Error
	return books, err
}

// GetBookByID 根据ID获取书籍（只返回上架状态）
// 参数:
//
//	id - 书籍ID
//
// 返回:
//
//	*model.Book - 书籍对象指针
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetBookByID(id int) (*model.Book, error) {
	var book model.Book
	// 对应SQL: SELECT * FROM books WHERE status = 1 AND id = id LIMIT 1;
	err := b.db.Where("status = ?", 1).First(&book, id).Error
	return &book, err
}

// GetBookByIDForAdmin 根据ID获取书籍（管理员用，不过滤状态）
// 参数:
//
//	id - 书籍ID
//
// 返回:
//
//	*model.Book - 书籍对象指针
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetBookByIDForAdmin(id int) (*model.Book, error) {
	var book model.Book
	// 对应SQL: SELECT * FROM books WHERE id = id LIMIT 1;
	err := b.db.First(&book, id).Error
	return &book, err
}

// GetBooksByType 根据类型获取书籍（只返回上架状态）
// 参数:
//
//	bookType - 书籍类型
//
// 返回:
//
//	[]*model.Book - 指定类型的书籍对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) GetBooksByType(bookType string) ([]*model.Book, error) {
	var books []*model.Book
	// 对应SQL: SELECT * FROM books WHERE status = 1 AND type = bookType;
	err := b.db.Where("status = ? AND type = ?", 1, bookType).Find(&books).Error
	return books, err
}

// SearchBooks 搜索书籍（只返回上架状态）
// 参数:
//
//	keyword - 搜索关键词
//
// 返回:
//
//	[]*model.Book - 符合条件的书籍对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) SearchBooks(keyword string) ([]*model.Book, error) {
	var books []*model.Book
	// 对应SQL: SELECT * FROM books WHERE status = 1 AND (title LIKE '%keyword%' OR author LIKE '%keyword%' OR description LIKE '%keyword%');
	err := b.db.Where("status = ? AND (title LIKE ? OR author LIKE ? OR description LIKE ?)", 1, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Find(&books).Error
	return books, err
}

// SearchBooksWithPagination 分页搜索书籍（只返回上架状态）
// 参数:
//
//	keyword - 搜索关键词
//	page - 页码，从1开始
//	pageSize - 每页记录数
//
// 返回:
//
//	[]*model.Book - 当前页的书籍对象切片
//	int64 - 符合条件的总记录数
//	error - 如果查询过程中出现错误则返回错误
func (b *BookDAO) SearchBooksWithPagination(keyword string, page, pageSize int) ([]*model.Book, int64, error) {
	var books []*model.Book
	var total int64

	// 构建搜索条件（只搜索上架状态）
	// 对应SQL条件: status = 1 AND (title LIKE '%keyword%' OR author LIKE '%keyword%' OR description LIKE '%keyword%')
	searchCondition := b.db.Where("status = ? AND (title LIKE ? OR author LIKE ? OR description LIKE ?)", 1, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	// 对应SQL: SELECT COUNT(*) FROM books WHERE status = 1 AND (title LIKE '%keyword%' OR author LIKE '%keyword%' OR description LIKE '%keyword%');
	err := searchCondition.Model(&model.Book{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	// 对应SQL: SELECT * FROM books WHERE status = 1 AND (title LIKE '%keyword%' OR author LIKE '%keyword%' OR description LIKE '%keyword%') LIMIT pageSize OFFSET offset;
	offset := (page - 1) * pageSize
	err = searchCondition.Offset(offset).Limit(pageSize).Find(&books).Error
	return books, total, err
}

// CreateBook 创建书籍
// 参数:
//
//	book - 书籍对象指针
//
// 返回:
//
//	error - 如果创建过程中出现错误则返回错误
func (b *BookDAO) CreateBook(book *model.Book) error {
	// 对应SQL: INSERT INTO books (title, author, price, ...) VALUES (book.Title, book.Author, book.Price, ...);
	err := b.db.Create(book).Error
	return err
}

// UpdateBook 更新书籍
// 参数:
//
//	book - 书籍对象指针
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (b *BookDAO) UpdateBook(book *model.Book) error {
	// 对应SQL: UPDATE books SET title = book.Title, author = book.Author, ... WHERE id = book.ID;
	err := b.db.Save(book).Error
	return err
}

// DeleteBook 删除书籍
// 参数:
//
//	id - 书籍ID
//
// 返回:
//
//	error - 如果删除过程中出现错误则返回错误
func (b *BookDAO) DeleteBook(id int) error {
	// 对应SQL: DELETE FROM books WHERE id = id;
	err := b.db.Delete(&model.Book{}, id).Error
	return err
}
