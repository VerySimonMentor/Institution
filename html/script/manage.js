const pageNum = 10;
var currentPage = 1;
var totalPage;

$(document).ready(function() {
    function fetchData(page = currentPage) {
        var data = {
            page: page,
            pageNum: pageNum
        };
        $.ajax({
            url: '/country/show',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify(data),
            success: function(data) {
                totalPage = data.totalPage;
                var table = $('#country-table tbody');
                table.empty();

                for (let i = 0; i < data.results.length; i++) {
                    var listIndex = (page - 1) * pageNum + i + 1;
                    var chiNameText = $(`<input type="text" class="input-text" value="${data.results[i].countryChiName}" readonly />`);
                    var engNameText = $(`<input type="text" class="input-text" value="${data.results[i].countryEngName}" readonly />`);
                    var row = $(
                        `<tr>
                            <td>${listIndex}</td>
                            <td>${chiNameText.prop('outerHTML')}</td>
                            <td>${engNameText.prop('outerHTML')}</td>
                            <td>${data.results[i].schoolNum}</td>
                            <td>${data.results[i].provinceNum}</td>
                            <td>
                                <a href="/country/edit/${data.results[i].id}" class="btn btn-primary">Edit</a>
                                <a href="/country/delete/${data.results[i].id}" class="btn btn-danger">Delete</a>
                            </td>
                        </tr>`
                    );
                    table.append(row);
                    row.find('input.input-text').eq(0).change((function(countryId, listIndex, countryName) {
                        return function() {
                            console.log(countryId, listIndex, countryName);
                            var value = $(this).val();
                            inputTextChange(countryId, listIndex, countryName, value);
                        }
                    })(data.results[i].countryId, listIndex, 'countryChiName'));
                    chiNameText.change((function(countryId, listIndex, countryName) {
                        return function() {
                            console.log(countryId, listIndex, countryName);
                            var value = $(this).val();
                            inputTextChange(countryId, listIndex, countryName, value);
                        }
                    })(data.results[i].countryId, listIndex, 'countryChiName'));
                    engNameText.change((function(countryId, listIndex, countryName, value) {
                        return function() {
                            inputTextChange(countryId, listIndex, countryName, value);
                        }
                    })(data.results[i].id, listIndex, 'countryEngName', engNameText.val()));

                    // Add event listener for text box change
                    $('.country-chi-name, .country-eng-name').on('change', function() {
                        // Handle text box change event here
                    });
                }

                $('.pagination').empty();
                for (let i = 1; i <= data.totalPage; i++) {
                    $('.pagination').append(`<a class="page-link" href="#" data-page="${i}">${i}</a></li>`);
                }
            }
        });

        $('#add-country-btn').prop('disabled', true);
        var countrySwitch = $('#country-switch input[type="checkbox"]');
        countrySwitch.prop('checked', false);
        countrySwitch.change(function() {
            var readonly = !$(this).prop('checked');
            $('#add-country-btn').prop('disabled', readonly);
            $('.input-text').prop('readonly', readonly);
        });
    }

    fetchData();

    $(document).on('click', '.page-link', function(e) {
        e.preventDefault();
        var page = $(this).data('page');
        fetchData(page);
    });

    $("#add-country-btn").click(function() {
        if (currentPage != totalPage) {
            fetchData(totalPage);
        }
        $.get("/country/create", function(response) {
            fetchData(totalPage);
        });
    });

    function inputTextChange(countryId, listIndex, field, value) {
        console.log(countryId, listIndex, field, value);
        $.ajax({
            url: '/country/update',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify({
                countryId: countryId,
                listIndex: listIndex,
                field: field,
                value: value
            }),
            success: function(data) {
                console.log(data);
            }
        });
    }
});