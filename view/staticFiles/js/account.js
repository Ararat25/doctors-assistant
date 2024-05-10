function check() {
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

function logout() {
    try {
        let param = new URL(location.href)
        let addr = new URL(param.origin + "/logout")
        addr.searchParams.set("user", param.searchParams.get("user"))
        fetch(addr, {method: 'Get'})
            .then(res => {
                if (res.ok) {
                    document.location.href = '/main'
                }
            })
            .catch(error => {
                console.log(error)
            })
    }
    catch (e) {
        console.log(e)
    }
}

// Примеры симптомов
const symptomExamples = [
    "Головная боль",
    "Температура",
    "Озноб",
    "Боль в горле",
    "Кашель",
    "Затрудненное дыхание",
    "Насморк",
    "Боль в животе",
    "Тошнота",
    "Рвота",
    "Диарея",
    "Боль в груди",
    "Повышенное сердцебиение",
    "Боль в спине",
    "Боль в мышцах и суставах",
    "Отеки",
    "Потеря аппетита",
    "Изменения веса",
    "Усталость",
    "Слабость",
    "Бессонница",
    "Изменения аппетита",
    "Покраснение кожи",
    "Сыпь",
    "Зуд",
    "Высыпания",
    "Жжение",
    "Боль при мочеиспускании",
    "Потеря сознания",
    "Повышенная чувствительность к свету"
];

symptomExamples.sort()

// Функция для отображения примеров симптомов
function displaySymptomExamples(symptoms) {
    const symptomContainer = document.getElementById("symptomExamples");
    symptomContainer.innerHTML = "";

    symptoms.forEach(symptom => {
        const symptomBlock = document.createElement("div");
        symptomBlock.classList.add("col", "mb-3");
        symptomBlock.innerHTML = `<div class="card">
                                <div class="card-body">
                                  <h5 class="card-title">${symptom}</h5>
                                </div>
                              </div>`;
        symptomBlock.addEventListener("click", function() {
            const inputField = document.getElementById("symptomInput");
            inputField.value += (inputField.value ? ", " : "") + (inputField.value ? symptom.toLowerCase() : symptom);
        });
        symptomBlock.addEventListener('mouseover', function() {
            const title = this.querySelector(".card-title");
            title.classList.add("expanded");
        });
        symptomBlock.addEventListener('mouseout', function() {
            const title = this.querySelector(".card-title");
            title.classList.remove("expanded");
        });
        symptomContainer.appendChild(symptomBlock);
    });
}


// Функция для фильтрации симптомов
function filterSymptoms(keyword) {
    const filteredSymptoms = symptomExamples.filter(symptom => {
        return symptom.toLowerCase().includes(keyword.toLowerCase());
    });
    displaySymptomExamples(filteredSymptoms);
}

// Обработчик события для поля фильтрации симптомов
document.getElementById("filterInput").addEventListener("input", function() {
    const keyword = this.value.trim();
    filterSymptoms(keyword);
});

// Начальное отображение примеров симптомов
displaySymptomExamples(symptomExamples);

function submitSymptoms() {
    const symptoms = document.getElementById("symptomInput").value;

    if (symptoms === "") {
        return
    }

    const requestData = {
        symptoms: symptoms
    };

    fetch('/account/assistant/message', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
        .then(response => response.json())
        .then(data => {
            let answer =  data.choices[0].Message.content
            answer = answer.split(/(?=\d\. )/)
            let htmlCode = '';

            answer.forEach(function(elem) {
                htmlCode += '<br>' + elem;
            });

            const diagnosisList = document.getElementById("diagnosisList");
            if (diagnosisList.textContent.trim() === "") {
                diagnosisList.innerHTML += "<b>Вы</b><br>" + symptoms + "<br><b>Ассистент</b>" + htmlCode
            }
            else {
                diagnosisList.innerHTML += "<br><br>" + "<b>Вы</b><br>" + symptoms + "<br><b>Ассистент</b>" + htmlCode
            }

        })
        .catch(error => console.error('Ошибка при отправке запроса:', error));
}