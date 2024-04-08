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
                var countrySwitch = $('#country-switch input[type="checkbox"]');
                var readonly = !countrySwitch.prop('checked');

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
                                <a href=# class="btn btn-province">编辑省份</a>
                                <a href=# class="btn btn-school">编辑学校</a>
                                <a href=# class="btn btn-delete">删除</a>
                            </td>
                        </tr>`
                    );
                    table.append(row);
                    row.find('input.input-text').eq(0).change((function(countryId, listIndex, countryName) {
                        return function() {
                            var value = $(this).val();
                            inputTextChange(countryId, listIndex, countryName, value);
                        }
                    })(data.results[i].countryId, listIndex, 'countryChiName'));
                    row.find('input.input-text').eq(1).change((function(countryId, listIndex, countryName) {
                        return function() {
                            var value = $(this).val();
                            inputTextChange(countryId, listIndex, countryName, value);
                        }
                    })(data.results[i].countryId, listIndex, 'countryEngName'));
                    row.find('input.input-text').prop('readonly', readonly);
                    row.find('.btn-delete').click((function(countryId, listIndex) {
                        return function() {
                            var data = {
                                countryId: countryId,
                                listIndex: listIndex - 1,
                            }
                            $.ajax({
                                url: '/country/delete',
                                type: 'DELETE',
                                headers: {
                                    'Content-Type': 'application/json'
                                },
                                data: JSON.stringify(data),
                                success: function(data) {
                                    fetchData(currentPage);
                                }
                            });
                        }
                    })(data.results[i].countryId, listIndex));
                }

                $('.pagination').empty();
                for (let i = 1; i <= data.totalPage; i++) {
                    $('.pagination').append(`<a class="page-link" href="#" data-page="${i}">${i}</a>`);
                }

                $('#add-country-btn').prop('disabled', readonly);
            }
        });
    }

    $(document).on('click', '.page-link', function(e) {
        e.preventDefault();
        currentPage = $(this).data('page');
        fetchData(currentPage);
    });

    $("#add-country-btn").click(function() {
        if (currentPage != totalPage) {
            currentPage = totalPage;
            fetchData(currentPage);
        }
        $.get("/country/create", {pageNum: pageNum}, function(response) {
            fetchData(response.totalPage);
        });
    });

    function inputTextChange(countryId, listIndex, field, value) {
        $.ajax({
            url: '/changeCountry',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify({
                countryId: countryId,
                listIndex: listIndex - 1,
                updateField: field,
                updateValue: value
            }),
            success: function(data) {
                // console.log(data);
            }
        });
    }

    function initCountry() {
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
    initCountry();
});