window.addEventListener('beforeunload', async function(event) {
    try {
        const response = await axios.get('/account');
        // обработка успешного ответа от сервера
    } catch (error) {
        if (error.response && error.response.status === 401) {
            // токен истек или недействителен
            event.preventDefault(); // отмена обновления страницы
            await refreshToken(); // обновление токена
            window.location.reload(); // повторная загрузка страницы
        }
    }
});

// Функция для обновления токена
async function refreshToken() {
    try {
        const response = await axios.get('/refresh-token');
    } catch (error) {
        // обработка ошибок при обновлении токена
        console.error('Ошибка при обновлении токена:', error);
    }
}