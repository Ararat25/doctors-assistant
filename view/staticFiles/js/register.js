async function handleFormSubmit(event) {
    event.preventDefault()
    const response = await fetch('register/user', {
        method: 'POST',
        body: new FormData(registerForm),
    })

    const usernameInput = document.getElementById('username');
    const errorElement = document.getElementById('status');

    if (response.status === 409) {
        errorElement.classList.add('text-danger');
        errorElement.textContent = 'Аккаунт с такой почтой уже существует';

        usernameInput.classList.add('is-invalid');

        usernameInput.addEventListener('input', () => {
            errorElement.textContent = "";
            usernameInput.classList.remove('is-invalid');
        });
        return null;
    }

    if (!response.ok) {
        throw new Error(`Ошибка по адресу 'register/user', статус ошибки ${response.status}`);
    }

    if (response.status === 200) {
        console.log("OK");
        registerForm.reset();
        return null
    }
}

function register() {
    const registerForm = document.getElementById('registerForm')
    registerForm.addEventListener('submit', handleFormSubmit)
}
