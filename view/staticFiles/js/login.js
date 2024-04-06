async function handleFormSubmit(event) {
    event.preventDefault()
    const response = await fetch('login/user', {
        method: 'POST',
        body: new FormData(loginForm),
    });

    if (!response.ok) {
        throw new Error(`Ошибка по адресу 'login/user', статус ошибки ${response.status}`);
    }
    loginForm.reset()
    return await response.json();
}

function login() {
    const loginForm = document.getElementById('loginForm')
    loginForm.addEventListener('submit', handleFormSubmit)
}
