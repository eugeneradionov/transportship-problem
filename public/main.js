suppliersCount = 1
consumersCount = 1
suppliersID = 1
consumersID = 1

$(document).ready(function () {
    renderShipmentMatrix()
})

$('#submit-button').on('click', function (event) {
    $.ajax({
        url: '/solve',
        method: 'POST',
        data: JSON.stringify(requestData()),
        success: function (data) {
            renderResultMatrix(data)
        },
    })
})

$(document).on('click', '[data-role="delete-supplier"]', function (event) {
    id = event.currentTarget.dataset.id
    $('[data-supplier-id=' + id + ']').remove()
    suppliersCount -= 1
    renderShipmentMatrix()
})

$(document).on('click', '[data-role="delete-consumer"]', function (event) {
    id = event.currentTarget.dataset.id
    $('[data-consumer-id=' + id + ']').remove()
    consumersCount -= 1
    renderShipmentMatrix()
})

$('[data-role="add-supplier"]').on('click', function (event) {
    suppliersCount += 1
    suppliersID += 1
    el = $('[data-role="supplier-placeholder"]').clone()
    el.removeClass('d-none')
    el.attr("data-role", null)
    el.attr("data-supplier-id", suppliersID)
    el.find('input').attr('data-role', 'supplier-stock')
    el.find('[data-role="delete-supplier"]').attr("data-id", suppliersID)
    el.find('[data-role="supplier-number"]').text(suppliersID)
    el.appendTo('[data-role="suppliers-list"]')
    renderShipmentMatrix()
})

$('[data-role="add-consumer"]').on('click', function (event) {
    consumersCount += 1
    consumersID += 1
    el = $('[data-role="consumer-placeholder"]').clone()
    el.removeClass('d-none')
    el.attr("data-role", null)
    el.attr("data-consumer-id", consumersID)
    el.find('input').attr('data-role', 'consumer-demand')
    el.find('[data-role="delete-consumer"]').attr("data-id", consumersID)
    el.find('[data-role="consumer-number"]').text(consumersID)
    el.appendTo('[data-role="consumers-list"]')
    renderShipmentMatrix()
})

function renderShipmentMatrix() {
    var table = $('<table></table>').addClass('table table-bordered table-striped')
    var header = $('<thead></thead>')

    header.append('<th scope="col">Suppliers/Consumers</th>')
    for (var i = 0; i < consumersCount; i++) {
        var headData = $('<th scope="col">' + (i+1) + '</th>')
        header.append(headData)
    }

    table.append(header)

    for (var i = 0; i < suppliersCount; i++) {
        row = $('<tr></tr>')
        row.append('<th scope="row">' + (i+1) + '</th>')
        for (var j = 0; j < consumersCount; j++) {
            var inputHTML = $('<input type="number" min="0" class="form-control" placeholder="Cost" aria-label="Cost">').attr('name', 'cost[' + i + '][]')
            var rowData = $('<td></td>').html(inputHTML)
            row.append(rowData)
        }
        table.append(row)
    }
    $('[data-role="shipment-matrix"]').html(table)
}

function renderResultMatrix(data) {
    var table = $('<table></table>').addClass('table table-bordered table-striped')
    var header = $('<thead></thead>')

    header.append('<th scope="col" class="text-center result-cell">From</th>')
    header.append('<th scope="col" class="text-center result-cell">To</th>')
    header.append('<th scope="col" class="text-center result-cell">Amount</th>')
    table.append(header)

    if (!data.transport_route) {
        notData = '<h3>No Data Provided</h3>'
        $('[data-role="modal-body"]').html(notData)
        return
    }

    data.transport_route.forEach(function (el) {
        row = $('<tr></tr>')
        var rowData = $('<td class="text-center result-cell"></td>').text(el.from.id)
        row.append(rowData)
        rowData = $('<td class="text-center result-cell"></td>').text(el.to.id)
        row.append(rowData)
        rowData = $('<td class="text-center result-cell"></td>').text(el.amount)
        row.append(rowData)
        table.append(row)
    })

    costRow = $('<tr></tr>')
    costRow.append('<th scope="row" class="text-center result-cell">Total Cost</th>')
    costRow.append('<td colspan="2" class="text-center result-cell"><strong>' + data.total_cost + '</strong></td>')
    table.append(costRow)

    $('[data-role="modal-body"]').html(table)
}

function requestData() {
    var supplierID = 1
    var consumerID = 1
    var suppliers = []
    var consumers = []
    var transportCost = []

    $('[data-role="supplier-stock"]').each(function (i, el) {
        suppliers.push({
            "id": supplierID,
            "stock": parseFloat(el.value)
        })
        supplierID += 1
    })

    $('[data-role="consumer-demand"]').each(function (i, el) {
        consumers.push({
            "id": consumerID,
            "demand": parseFloat(el.value)
        })
        consumerID += 1
    })

    for (let i = 0; i < suppliersCount; i++) {
        transportCost[i] = []
        let name = 'cost[' + i + '][]'
        $('input[name="' + name + '"]').each(function (j, el) {
            transportCost[i][j] = parseFloat(el.value)
        })
    }

    return {
        'suppliers': suppliers,
        'consumers': consumers,
        'transport_cost': transportCost
    }
}
