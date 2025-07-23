document.addEventListener('DOMContentLoaded', function() {
    // Загрузка списка городов
    loadCities();
    
    // Обработка формы
    const form = document.getElementById('delivery-form');
    form.addEventListener('submit', function(e) {
        e.preventDefault();
        calculateDelivery();
    });
});
async function loadCities(retryCount = 3) {
    const select = document.getElementById('city');
    
    try {
        select.innerHTML = '<option value="" disabled selected>Загрузка городов...</option>';
        
        let response;
        try {
            response = await fetch('/api/cities');
        } catch (error) {
            if (retryCount > 0) {
                return loadCities(retryCount - 1);
            }
            throw error;
        }

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || 'Ошибка сервера');
        }

        const cities = await response.json();
        
        select.innerHTML = '';
        const defaultOption = new Option('Выберите город', '', true, true);
        defaultOption.disabled = true;
        select.add(defaultOption);
        
        cities.forEach(city => {
            select.add(new Option(city, city));
        });

        // Автовыбор Москвы
        const moscowOption = Array.from(select.options)
            .find(opt => opt.value === 'Москва');
        if (moscowOption) moscowOption.selected = true;

    } catch (error) {
        console.error('Ошибка загрузки городов:', error);
        select.innerHTML = `
            <option value="" disabled selected>
                Ошибка: ${error.message}
            </option>
        `;
        select.add(new Option('Москва', 'Москва'));
    }
}
async function calculateDelivery() {
    const city = document.getElementById('city').value;
    const weight = document.getElementById('weight').value;
    const resultDiv = document.getElementById('result');
    
    // Валидация
    if (!city || !weight) {
        resultDiv.innerHTML = '<p class="error">Пожалуйста, заполните все поля</p>';
        return;
    }

    // Показываем индикатор загрузки
    resultDiv.innerHTML = '<p>Идет расчет...</p>';
    resultDiv.className = 'loading';

    try {
        const response = await fetch(`/api/calculate?city=${encodeURIComponent(city)}&weight=${weight}`);
        
        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || `Ошибка сервера: ${response.status}`);
        }

        const data = await response.json();
        
        // Отображаем результат
        resultDiv.innerHTML = `
            <p><strong>Результат:</strong> ${data.message}</p>
            <p><strong>Стоимость:</strong> ${data.price} руб.</p>
        `;
        resultDiv.className = data.status === 'OK' ? 'success' : 'error';

    } catch (error) {
        console.error('Ошибка расчета:', error);
        resultDiv.innerHTML = `<p class="error">Ошибка: ${error.message}</p>`;
        resultDiv.className = 'error';
        
        // Для ошибок API предлагаем повторить
        if (error.message.includes('Ошибка расчета стоимости доставки')) {
            resultDiv.innerHTML += '<button onclick="calculateDelivery()">Повторить</button>';
        }
    }
}