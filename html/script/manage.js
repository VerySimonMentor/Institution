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
                    table.append(
                        `<tr>
                            <td>${(page - 1) * pageNum + i + 1}</td>
                            <td><input type="text" class="input-text" value="${data.results[i].countryChiName}" data-field="countryChiName" /></td>
                            <td><input type="text" class="input-text" value="${data.results[i].countryEngName}" /></td>
                            <td>${data.results[i].schoolNum}</td>
                            <td>${data.results[i].provinceNum}</td>
                            <td>
                                <a href="/country/edit/${data.results[i].id}" class="btn btn-primary">Edit</a>
                                <a href="/country/delete/${data.results[i].id}" class="btn btn-danger">Delete</a>
                            </td>
                        </tr>`
                    );

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

    $(document).on('change', '.input-text', function() {
        console.log('change');
    });
});