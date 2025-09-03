function generateOrders() {

    var genOutput = document.getElementById('generateOrderMsg');

    fetch(`/generate`, {
        method: "POST"
    })
        .then(response => response.json())
        .then(response => {
            if (response.error) {
                console.error('Ошибка:', error);
            }
            else {
                console.log(`Сгенерировано 10 заказов`);
                genOutput.value = `Сгенерировано 10 заказов`;
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}

function getOrder() {

    var orderIdInput = document.getElementById('orderUid').value;

    if (orderIdInput.trim() == "") {
        alert("Неверный номер заказа");
        return;
    }

    var orderUid = String(orderIdInput);

    fetch(`/orders/${orderUid}`, {
        method: "GET"
    })
        .then(response => response.json())
        .then(response => {
            if (response.error) {
                console.error('Ошибка:', response.message);
            }
            else {
                var order = response.message;

                clearTables();

                fillOrderInfo(order);
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}

function fillOrderInfo(order) {

    // Заказы

    var ignoredFields = ['delivery', 'payment', 'items'];
    fillObject('orderTable', ignoredFields, order);

    // Доставка

    fillObject('deliveryTable', [], order.delivery);

    // Оплата

    fillObject('paymentTable', [], order.payment);

    // Товары

    const itemsBody = document.getElementById('itemsTable').querySelector('tbody');
    order.items.forEach((prop) => {

        const row = document.createElement('tr');
        row.insertCell().textContent = prop.chrt_id;
        row.insertCell().textContent = prop.track_number;
        row.insertCell().textContent = prop.price;
        row.insertCell().textContent = prop.rid;
        row.insertCell().textContent = prop.name;
        row.insertCell().textContent = prop.sale;
        row.insertCell().textContent = prop.size;
        row.insertCell().textContent = prop.total_price;
        row.insertCell().textContent = prop.nm_id;
        row.insertCell().textContent = prop.brand;
        row.insertCell().textContent = prop.status;
        itemsBody.appendChild(row);
    });
}

function fillObject(tableNmae, ignoredFields, entry) {
    var table = document.getElementById(tableNmae).querySelector('tbody');
    row = document.createElement('tr');
    Object.entries(entry).forEach(([key, value]) => {

        if (!ignoredFields.includes(key)) {
            const cell = document.createElement('td');
            cell.textContent = value;
            row.appendChild(cell);
        }
    });
    table.appendChild(row);
}

function clearTables() {
    document.getElementById('orderTable').querySelector('tbody').innerHTML = '';
    document.getElementById('deliveryTable').querySelector('tbody').innerHTML = '';
    document.getElementById('paymentTable').querySelector('tbody').innerHTML = '';
    document.getElementById('itemsTable').querySelector('tbody').innerHTML = '';
}