async function check() {
    try {
        let param = new URL(location.href)
        let addr = new URL(param.origin + "/account/user")
        addr.searchParams.set("user", param.searchParams.get("user"))
        fetch(addr, { method: 'GET',})
            .then((res) => {
                if (res.status === 401) {
                    fetch('/refresh-token', {
                        method: 'GET',
                    }).then((respRefreshToken) => {
                        if (!respRefreshToken.ok) {
                            respRefreshToken.text().then(text => {
                                alert("Время вашего сеанса истекло")
                                window.location.href = '/login';
                            })
                        }
                    })
                }
                else if (!res.ok) {
                    res.text().then(text => {
                        throw new Error(`Server error status: ${res.status} message: ${text}`)
                    })
                }
            })
            .catch(error => {
                console.log('Error:', error);
            })
    }
    catch (e) {
        console.log(e)
    }
}
setInterval(check, 1000*59)