package service

import (
	"bookstore/model"
	"bookstore/repository"
)

// BookService 书籍服务
// 封装了所有与书籍相关的业务逻辑操作
type BookService struct {
	BookDB     *repository.BookDAO     // 书籍数据访问对象
	CategoryDB *repository.CategoryDAO // 分类数据访问对象
}

// NewBookService 创建新的书籍服务实例
// 返回:
//
//	*BookService - 初始化后的书籍服务对象
func NewBookService() *BookService {
	return &BookService{
		BookDB:     repository.NewBookDAO(),
		CategoryDB: repository.NewCategoryDAO(),
	}
}

// GetAllBooks 获取所有书籍
// 返回:
//
//	[]*model.Book - 书籍对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) GetAllBooks() ([]*model.Book, error) {
	return b.BookDB.GetAllBooks()
}

// GetBooksByPage 分页获取书籍
// 参数:
//
//	page - 页码
//	pageSize - 每页大小
//
// 返回:
//
//	[]*model.Book - 书籍对象切片
//	int64 - 书籍总数
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) GetBooksByPage(page, pageSize int) ([]*model.Book, int64, error) {
	return b.BookDB.GetBooksByPage(page, pageSize)
}

// GetHotBooks 获取热销书籍
// 参数:
//
//	limit - 限制数量
//
// 返回:
//
//	[]*model.Book - 热销书籍对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) GetHotBooks(limit int) ([]*model.Book, error) {
	return b.BookDB.GetHotBooks(limit)
}

// GetNewBooks 获取新书
// 参数:
//
//	limit - 限制数量
//
// 返回:
//
//	[]*model.Book - 新书对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) GetNewBooks(limit int) ([]*model.Book, error) {
	return b.BookDB.GetNewBooks(limit)
}

// GetBookByID 根据ID获取书籍
// 参数:
//
//	id - 书籍ID
//
// 返回:
//
//	*model.Book - 书籍对象指针
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) GetBookByID(id int) (*model.Book, error) {
	return b.BookDB.GetBookByID(id)
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
func (b *BookService) GetBookByIDForAdmin(id int) (*model.Book, error) {
	return b.BookDB.GetBookByIDForAdmin(id)
}

// GetBooksByType 根据类型获取书籍
// 参数:
//
//	bookType - 书籍类型
//
// 返回:
//
//	[]*model.Book - 书籍对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) GetBooksByType(bookType string) ([]*model.Book, error) {
	return b.BookDB.GetBooksByType(bookType)
}

// SearchBooks 搜索书籍
// 参数:
//
//	keyword - 搜索关键词
//
// 返回:
//
//	[]*model.Book - 书籍对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) SearchBooks(keyword string) ([]*model.Book, error) {
	return b.BookDB.SearchBooks(keyword)
}

// SearchBooksWithPagination 分页搜索书籍
// 参数:
//
//	keyword - 搜索关键词
//	page - 页码
//	pageSize - 每页大小
//
// 返回:
//
//	[]*model.Book - 书籍对象切片
//	int64 - 搜索结果总数
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) SearchBooksWithPagination(keyword string, page, pageSize int) ([]*model.Book, int64, error) {
	return b.BookDB.SearchBooksWithPagination(keyword, page, pageSize)
}

// CreateBook 创建书籍
// 参数:
//
//	book - 书籍对象指针
//
// 返回:
//
//	error - 如果创建过程中出现错误则返回错误
func (b *BookService) CreateBook(book *model.Book) error {
	return b.BookDB.CreateBook(book)
}

// UpdateBook 更新书籍
// 参数:
//
//	book - 书籍对象指针
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (b *BookService) UpdateBook(book *model.Book) error {
	return b.BookDB.UpdateBook(book)
}

// DeleteBook 删除书籍
// 参数:
//
//	id - 书籍ID
//
// 返回:
//
//	error - 如果删除过程中出现错误则返回错误
func (b *BookService) DeleteBook(id int) error {
	return b.BookDB.DeleteBook(id)
}

// CreateBookFromRequest 从请求创建图书
// 参数:
//
//	req - 图书创建请求对象指针
//
// 返回:
//
//	error - 如果创建过程中出现错误则返回错误
func (b *BookService) CreateBookFromRequest(req *model.BookCreateRequest) error {
	// 确保CategoryID有有效值
	categoryID := req.CategoryID
	if categoryID <= 0 {
		categoryID = 1 // 默认使用第一个分类
	}

	book := &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Price:       req.Price,
		Discount:    req.Discount,
		Type:        req.Type,
		Stock:       req.Stock,
		Status:      req.Status,
		CoverURL:    req.CoverURL,
		Description: req.Description,
		ISBN:        req.ISBN,
		Publisher:   req.Publisher,
		PublishDate: req.PublishDate,
		Pages:       req.Pages,
		Language:    req.Language,
		Format:      req.Format,
		CategoryID:  categoryID,
		Sale:        req.Sale,
	}

	return b.BookDB.CreateBook(book)
}

// UpdateBookFromRequest 从请求更新图书
// 参数:
//
//	id - 书籍ID
//	req - 图书更新请求对象指针
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (b *BookService) UpdateBookFromRequest(id uint, req *model.BookUpdateRequest) error {
	book, err := b.BookDB.GetBookByIDForAdmin(int(id))
	if err != nil {
		return err
	}

	// 更新字段（排除状态字段，避免意外修改状态）
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Price > 0 {
		book.Price = req.Price
	}
	if req.Discount > 0 {
		book.Discount = req.Discount
	}
	if req.Type != "" {
		book.Type = req.Type
	}
	if req.Stock >= 0 {
		book.Stock = req.Stock
	}
	// 注意：不更新 Status 字段，避免意外修改状态
	// if req.Status >= 0 {
	// 	book.Status = req.Status
	// }
	if req.CoverURL != "" {
		book.CoverURL = req.CoverURL
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.ISBN != "" {
		book.ISBN = req.ISBN
	}
	if req.Publisher != "" {
		book.Publisher = req.Publisher
	}
	if req.PublishDate != "" {
		book.PublishDate = req.PublishDate
	}
	if req.Pages > 0 {
		book.Pages = req.Pages
	}
	if req.Language != "" {
		book.Language = req.Language
	}
	if req.Format != "" {
		book.Format = req.Format
	}
	if req.CategoryID > 0 {
		book.CategoryID = req.CategoryID
	}
	if req.Sale >= 0 {
		book.Sale = req.Sale
	}

	return b.BookDB.UpdateBook(book)
}

// GetBookList 获取图书列表（管理员）
// 参数:
//
//	req - 图书列表请求对象指针
//
// 返回:
//
//	*model.BookListResponse - 图书列表响应对象指针
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) GetBookList(req *model.BookListRequest) (*model.BookListResponse, error) {
	books, total, err := b.BookDB.GetBooksByPageForAdminWithSearch(req.Page, req.PageSize, req.Title, req.Author, req.Type, req.Status)
	if err != nil {
		return nil, err
	}

	// 转换为Book切片
	var bookList []model.Book
	for _, book := range books {
		bookList = append(bookList, *book)
	}

	totalPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &model.BookListResponse{
		Books:       bookList,
		Total:       total,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}, nil
}

// UpdateBookStatus 更新图书状态
// 参数:
//
//	id - 书籍ID
//	status - 状态值
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (b *BookService) UpdateBookStatus(id uint, status int) error {
	book, err := b.BookDB.GetBookByIDForAdmin(int(id))
	if err != nil {
		return err
	}

	book.Status = status
	return b.BookDB.UpdateBook(book)
}

// GetCategories 获取所有分类
// 返回:
//
//	[]model.Category - 分类对象切片
//	error - 如果查询过程中出现错误则返回错误
func (b *BookService) GetCategories() ([]model.Category, error) {
	categories, err := b.CategoryDB.GetAllCategories()
	if err != nil {
		return nil, err
	}

	// 转换为Category切片
	var categoryList []model.Category
	for _, category := range categories {
		categoryList = append(categoryList, *category)
	}

	return categoryList, nil
}

// CreateCategory 创建分类
// 参数:
//
//	category - 分类对象指针
//
// 返回:
//
//	error - 如果创建过程中出现错误则返回错误
func (b *BookService) CreateCategory(category *model.Category) error {
	return b.CategoryDB.CreateCategory(category)
}

// UpdateCategory 更新分类
// 参数:
//
//	id - 分类ID
//	updates - 更新字段映射
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (b *BookService) UpdateCategory(id uint, updates map[string]any) error {
	category, err := b.CategoryDB.GetCategoryByID(int(id))
	if err != nil {
		return err
	}

	// 更新字段
	if name, ok := updates["name"].(string); ok && name != "" {
		category.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		category.Description = description
	}
	if icon, ok := updates["icon"].(string); ok {
		category.Icon = icon
	}
	if color, ok := updates["color"].(string); ok {
		category.Color = color
	}
	if gradient, ok := updates["gradient"].(string); ok {
		category.Gradient = gradient
	}
	if sort, ok := updates["sort"].(int); ok {
		category.Sort = sort
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		category.IsActive = isActive
	}

	return b.CategoryDB.UpdateCategory(category)
}

// DeleteCategory 删除分类
// 参数:
//
//	id - 分类ID
//
// 返回:
//
//	error - 如果删除过程中出现错误则返回错误
func (b *BookService) DeleteCategory(id uint) error {
	return b.CategoryDB.DeleteCategory(int(id))
}
