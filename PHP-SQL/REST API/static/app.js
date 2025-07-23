const API_BASE_URL = 'http://localhost:8080';

document.addEventListener('DOMContentLoaded', function() {
    const productForm = document.getElementById('product-form');
    const productsList = document.getElementById('products-list');
    const loadingIndicator = document.getElementById('loading');
    const submitBtn = document.getElementById('submit-btn');
    const cancelBtn = document.getElementById('cancel-btn');
    const formTitle = document.getElementById('form-title');
    
    let isEditing = false;
    let currentProductId = null;
    
    // Загрузка товаров при загрузке страницы
    loadProducts();
    
    // Обработка отправки формы
    productForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const product = {
            name: document.getElementById('name').value,
            price: parseFloat(document.getElementById('price').value),
            stock: parseInt(document.getElementById('stock').value)
        };
        
        if (isEditing) {
            updateProduct(currentProductId, product);
        } else {
            createProduct(product);
        }
    });
    
    // Обработка отмены редактирования
    cancelBtn.addEventListener('click', resetForm);
    
    // Функция загрузки товаров
    function loadProducts() {
        loadingIndicator.style.display = 'block';
        productsList.innerHTML = '';
        
        fetch(`${API_BASE_URL}/api/products`) // Обновите этот URL
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка загрузки товаров');
                }
                return response.json();
            })
            .then(products => {
                loadingIndicator.style.display = 'none';
                
                if (products.length === 0) {
                    productsList.innerHTML = '<tr><td colspan="5">Нет товаров</td></tr>';
                    return;
                }
                
                products.forEach(product => {
                    const row = document.createElement('tr');
                    
                    row.innerHTML = `
                        <td>${product.id}</td>
                        <td>${product.name}</td>
                        <td>${product.price.toFixed(2)}</td>
                        <td>${product.stock}</td>
                        <td>
                            <button onclick="editProduct(${product.id})" class="edit-btn">Редактировать</button>
                            <button onclick="deleteProduct(${product.id})" class="danger">Удалить</button>
                        </td>
                    `;
                    
                    productsList.appendChild(row);
                });
            })
            .catch(error => {
                loadingIndicator.style.display = 'none';
                productsList.innerHTML = `<tr><td colspan="5" class="error-message">${error.message}</td></tr>`;
                console.error('Ошибка:', error);
            });
    }
    
    // Функция создания товара
    function createProduct(product) {
        fetch(`${API_BASE_URL}/api/products`, { // Обновите этот URL
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(product)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка создания товара');
            }
            return response.json();
        })
        .then(() => {
            showMessage('Товар успешно создан', true);
            resetForm();
            loadProducts();
        })
        .catch(error => {
            showMessage(error.message, false);
            console.error('Ошибка:', error);
        });
    }
    
    // Обновите аналогично все остальные fetch-запросы:
    // updateProduct, deleteProduct, editProduct
    function updateProduct(id, product) {
        fetch(`${API_BASE_URL}/api/products/${id}`, { // Обновите
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(product)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка обновления товара');
            }
            return response.json();
        })
        .then(() => {
            showMessage('Товар успешно обновлен', true);
            resetForm();
            loadProducts();
        })
        .catch(error => {
            showMessage(error.message, false);
            console.error('Ошибка:', error);
        });
    }
    
    window.deleteProduct = function(id) {
        if (!confirm('Вы уверены, что хотите удалить этот товар?')) {
            return;
        }
        
        fetch(`${API_BASE_URL}/api/products/${id}`, { // Обновите
            method: 'DELETE'
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка удаления товара');
            }
            showMessage('Товар успешно удален', true);
            loadProducts();
        })
        .catch(error => {
            showMessage(error.message, false);
            console.error('Ошибка:', error);
        });
    };
    
    window.editProduct = function(id) {
        fetch(`${API_BASE_URL}/api/products/${id}`) // Обновите
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка загрузки товара');
                }
                return response.json();
            })
            .then(product => {
                isEditing = true;
                currentProductId = product.id;
                
                document.getElementById('product-id').value = product.id;
                document.getElementById('name').value = product.name;
                document.getElementById('price').value = product.price;
                document.getElementById('stock').value = product.stock;
                
                formTitle.textContent = 'Редактировать товар';
                submitBtn.textContent = 'Обновить';
                cancelBtn.style.display = 'inline-block';
                
                document.getElementById('product-form').scrollIntoView({ behavior: 'smooth' });
            })
            .catch(error => {
                showMessage(error.message, false);
                console.error('Ошибка:', error);
            });
    };
    
    // Сброс формы
    function resetForm() {
        productForm.reset();
        isEditing = false;
        currentProductId = null;
        
        formTitle.textContent = 'Добавить новый товар';
        submitBtn.textContent = 'Сохранить';
        cancelBtn.style.display = 'none';
    }
    
    // Показать сообщение
    function showMessage(message, isSuccess) {
        const messageDiv = document.createElement('div');
        messageDiv.className = isSuccess ? 'success-message' : 'error-message';
        messageDiv.textContent = message;
        
        const formActions = document.querySelector('.form-actions');
        formActions.appendChild(messageDiv);
        
        // Удалить сообщение через 3 секунды
        setTimeout(() => {
            messageDiv.remove();
        }, 3000);
    }
});