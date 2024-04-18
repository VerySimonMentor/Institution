const pageNum = 10;

var currentCountryPage = 1;
var currentSchoolPage = 1;
var currentItemPage = 1;
var currentUserPage = 1;
var totalCountryPage;
var totalSchoolPage;
var totalItemPage;
var totalUserPage;

$(document).ready(function() {
    function fetchCountryData(page = currentCountryPage) {
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
                totalCountryPage = data.totalPage;
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
                            countryTextChange(countryId, listIndex, countryName, value);
                        }
                    })(data.results[i].countryId, listIndex, 'countryChiName'));
                    row.find('input.input-text').eq(1).change((function(countryId, listIndex, countryName) {
                        return function() {
                            var value = $(this).val();
                            countryTextChange(countryId, listIndex, countryName, value);
                        }
                    })(data.results[i].countryId, listIndex, 'countryEngName'));
                    row.find('input.input-text').prop('readonly', readonly);
                    row.find('.btn-delete').click((function(countryId, listIndex) {
                        return function() {
                            alert('确定删除吗？');
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
                                    if (table.children().length === 1) {
                                        $('#country-pagination').children().last().remove();
                                        if (currentCountryPage > 1) {
                                            currentCountryPage--;
                                        }
                                    }
                                    fetchCountryData(currentCountryPage);
                                }
                            });
                        }
                    })(data.results[i].countryId, listIndex));
                    row.find('.btn-school').click((function(listIndex) {
                        return function() {
                            $("#manage-country-content").css("display", "none");
                            $("#manage-school-content").css("display", "block");
                            $("#manage-school-content").css({
                                "position": "absolute", // 使用绝对定位
                                "top": "80px", 
                                "left": "20vw" 
                            });
                            initSchool(listIndex);
                        }
                    })(listIndex));
                    row.find('.btn-province').off('click').click((function(countryId, listIndex, currentPage) {
                        return function() {
                            disableSidebarButton();
                            var data = {
                                countryId: countryId,
                                listIndex: listIndex - 1,
                            }
                            $.ajax({
                                url: '/country/changeProvince/show',
                                type: 'POST',
                                headers: {
                                    'Content-Type': 'application/json'
                                },
                                data: JSON.stringify(data),
                                success: function(data) {
                                    $('#province-model').css('display', 'block');
                                    $("#manage-country-content").css('pointer-events', 'none');
                                    $("#manage-country-content").css('opacity', '0.5');
                                    $('#chinese-name-input').val(data.country.countryChiName);
                                    $('#english-name-input').val(data.country.countryEngName);
                                    $('#save-province-btn').off('click').click((function(countryId, listIndex){
                                        return function(){
                                            enableSidebarButton();
                                            var province = getProvinceData();
                                            var countryChiName = $('#chinese-name-input').val();
                                            var countryEngName = $('#english-name-input').val();
                                            var data = {
                                                countryId: countryId,
                                                listIndex: listIndex-1,
                                                countryChiName: countryChiName,
                                                countryEngName: countryEngName,
                                                province: province
                                            }
                                            $.ajax({
                                                url: '/country/changeProvince/save',
                                                type: 'POST',
                                                headers: {
                                                    'Content-Type': 'application/json'
                                                },
                                                data: JSON.stringify(data),
                                                success: function(data) {
                                                    $('#province-model').css('display', 'none');
                                                    $("#manage-country-content").css('pointer-events', 'auto');
                                                    $("#manage-country-content").css('opacity', '1');
                                                    fetchCountryData(currentPage);
                                                }
                                            });
                                        }
                                    })(data.country.countryId, listIndex));
                                    fetchProvinceData(data.country.province);
                                }
                            });
                        }//
                    })(data.results[i].countryId, listIndex, currentCountryPage));
                }

                $('#country-pagination').empty();
                for (let i = 1; i <= data.totalPage; i++) {
                    $('#country-pagination').append(`<a class="page-link" href="#" data-page="${i}">${i}</a>`);
                }

                $('#add-country-btn').prop('disabled', readonly);
            }
        });
    }

    $(document).on('click', '#country-pagination .page-link', function(e) {
        e.preventDefault();
        currentCountryPage = $(this).data('page');
        fetchCountryData(currentCountryPage);
    });

    $("#add-country-btn").click(function() {
        if (currentCountryPage != totalCountryPage) {
            currentCountryPage = totalCountryPage;
            fetchCountryData(currentCountryPage);
        }
        $.get("/country/create", {pageNum: pageNum}, function(response) {
            currentCountryPage = response.totalPage;
            fetchCountryData(response.totalPage);
        });
    });

    function countryTextChange(countryId, listIndex, field, value) {
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

    $(".button-list-button").click(function() {
        //侧边栏按钮section跳转
        //获取button id
        var buttonId = $(this).attr("id");
        //隐藏所有section
        $("#manage-country-content").css("display", "none");
        $("#manage-school-content").css("display", "none");
        $("#manage-item-content").css("display", "none");
        $("#manage-user-content").css("display", "none");
        $("#system-set-content").css("display", "none");
        //根据button id显示对应的section
        $('#' + buttonId + '-content') 
            .css("display", "block")
            .css({
                "position": "absolute", // 使用绝对定位
                "top": "80px", 
                "left": "20vw",
            });
        switch (buttonId) {
            case 'manage-country':
                fetchCountryData();
                initCountry();
                break;
            case 'manage-school':
                initSchool();
                break;
            case 'manage-item':
                initItem();
                break;
            case 'manage-user':
                initUser();
                break;
            case 'system-set':
                initSystemSet();
                break;
            default:
                break;
        }
    });

    function fetchProvinceData(province) {
        var table = $('#province-table tbody');
        table.empty();
        for (let i=0; i<province.length; i++) {
            var chiNameText = $(`<input type="text" class="input-text" value="${province[i].chiName}" />`);
            var engNameText = $(`<input type="text" class="input-text" value="${province[i].engName}" />`);
            var row = $(
                `<tr>
                    <td>${i+1}</td>
                    <td>${chiNameText.prop('outerHTML')}</td>
                    <td>${engNameText.prop('outerHTML')}</td>
                    <td>
                        <a href=# class="btn btn-province-delete">删除</a>
                    </td>
                </tr>`
            );
            row.find('.btn-province-delete').off('click').click((function(i){
                return function() {
                    var province = getProvinceData();
                    province.splice(i, 1);
                    fetchProvinceData(province);
                }
            })(i));
            table.append(row);
        }
    }

    function getProvinceData() {
        var province = [];
        $('#province-table tbody tr').each(function() {
            var chiName = $(this).find('input.input-text').eq(0).val();
            var engName = $(this).find('input.input-text').eq(1).val();
            province.push({chiName: chiName, engName: engName});
        });
        return province;
    }

    $('#add-province-btn').click(function(){
        var province = getProvinceData();
        province.push({chiName: "新省份", engName: "New Province"});
        fetchProvinceData(province);
    })

    $('#cancel-province-btn').click(function(){
        enableSidebarButton();
        $('#province-model').css('display', 'none');
        $("#manage-country-content").css('pointer-events', 'auto');
        $("#manage-country-content").css('opacity', '1');
    })

    function initSchool(listIndex = 0) {
        $('#school-table').css('width', '1600px');
        $.ajax({
            url: '/school/initPage',
            type: 'GET',
            success: function(data) {
                $('#add-school-btn').prop('disabled', true);
                var schoolSwitch = $('#school-switch input[type="checkbox"]');
                schoolSwitch.prop('checked', false);
                schoolSwitch.change(function() {
                    var readonly = !$(this).prop('checked');
                    $('#add-school-btn').prop('disabled', readonly);
                    $('.input-text').prop('readonly', readonly);
                    $('.input-select').prop('disabled', readonly);
                });
                var allCountry = data.results;
                var countrySelect = $('#school-page-country-select');
                countrySelect.empty();
                countrySelect.append('<option value="0">请选择国家</option>');
                for (var i = 0; i < allCountry.length; i++) {
                    var option = $(`<option value="${i+1}">${allCountry[i]}</option>`);
                    countrySelect.append(option);
                }
                countrySelect.off('change').change(function() {
                    var listIndex = $(this).val();
                    if (listIndex > 0) {
                        fetchSchoolData(listIndex);
                    } else {
                        $('#school-table tbody').empty();
                    }
                });
                countrySelect.val(listIndex).trigger('change');
            }
        });
    }

    function fetchSchoolData(listIndex, page = 1) {
        var data = {
            countryListIndex: listIndex - 1,
            page: page,
            pageNum: pageNum
        }
        $.ajax({
            url: '/country/editSchool',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify(data),
            success: function(data) {
                var school = data.results;
                var province = data.province;
                var schoolTypeList = data.schoolTypeList;
                totalSchoolPage = data.totalPage;
                var table = $('#school-table tbody');
                table.empty();
                var schoolSwitch = $('#school-switch input[type="checkbox"]');
                var readonly = !schoolSwitch.prop('checked');

                for (let i=0; i<school.length; i++) {
                    var listIndex = (page - 1) * pageNum + i + 1;
                    var chiNameText = $(`<input type="text" class="input-text" value="${school[i].schoolChiName}" />`);
                    var engNameText = $(`<input type="text" class="input-text" value="${school[i].schoolEngName}" />`);
                    var abbreviationText = $(`<input type="text" class="input-text" value="${school[i].schoolAbbreviation}" />`);
                    var typeSelect = $(`<select class="input-select"></select>`);
                    typeSelect.append('<option value=-1>请选择学校类型</option>');
                    for (var j = 0; j < schoolTypeList.length; j++) {
                        var option = $(`<option value="${schoolTypeList[j].schoolTypeId}" ${school[i].schoolType == schoolTypeList[j].schoolTypeId ? 'selected' : ''}>${schoolTypeList[j].schoolTypeName}</option>`);
                        typeSelect.append(option);
                    }
                    var provinceSelect = $(`<select class="input-select"></select>`);
                    for (var j = 0; j < province.length; j++) {
                        var option = $(`<option value="${province[j].chiName}" ${school[i].province === province[j].chiName ? 'selected' : ''}>${province[j].chiName}</option>`);
                        provinceSelect.append(option);
                    }
                    var linkText = $(`<input type="text" class="input-text" value="${school[i].officialWebLink}" />`);
                    var remarkText = $(`<input type="text" class="input-text" value="${school[i].schoolRemark}" />`);
                    
                    var row = $(
                        `<tr>
                            <td>${i+1+(page-1)*pageNum}</td>
                            <td>${chiNameText.prop('outerHTML')}</td>
                            <td>${engNameText.prop('outerHTML')}</td>
                            <td>${abbreviationText.prop('outerHTML')}</td>
                            <td>${typeSelect.prop('outerHTML')}</td>
                            <td>${provinceSelect.prop('outerHTML')}</td>
                            <td>${linkText.prop('outerHTML')}</td>
                            <td>${remarkText.prop('outerHTML')}</td>
                            <td>${school[i].itemNum}</td>
                            <td>
                                <a href=# class="btn btn-item">编辑项目</a>
                                <a href=# class="btn btn-school-delete">删除</a>
                            </td>
                        </tr>`
                    );
                    table.append(row);
                    row.find('input.input-text').eq(0).change((function(schoolId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            schoolTextChange(schoolId, listIndex, field, value);
                        }
                    })(school[i].schoolId, listIndex, 'schoolChiName'));
                    row.find('input.input-text').eq(1).change((function(schoolId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            schoolTextChange(schoolId, listIndex, field, value);
                        }
                    })(school[i].schoolId, listIndex, 'schoolEngName'));
                    row.find('input.input-text').eq(2).change((function(schoolId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            schoolTextChange(schoolId, listIndex, field, value);
                        }
                    })(school[i].schoolId, listIndex, 'schoolAbbreviation'));
                    row.find('select.input-select').eq(0).change((function(schoolId, listIndex, field) {
                        return function() {
                            var value = parseInt($(this).val());
                            if (value >= 0) {
                                schoolTextChange(schoolId, listIndex, field, value);
                            }
                        }
                    })(school[i].schoolId, listIndex, 'schoolType'));
                    row.find('select.input-select').eq(1).change((function(schoolId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            schoolTextChange(schoolId, listIndex, field, value);
                        }
                    })(school[i].schoolId, listIndex, 'province'));
                    row.find('input.input-text').eq(3).change((function(schoolId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            schoolTextChange(schoolId, listIndex, field, value);
                        }
                    })(school[i].schoolId, listIndex, 'officialWebLink'));
                    row.find('input.input-text').eq(4).change((function(schoolId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            schoolTextChange(schoolId, listIndex, field, value);
                        }
                    })(school[i].schoolId, listIndex, 'schoolRemark'));
                    row.find('input.input-text').prop('readonly', readonly);
                    row.find('select.input-select').prop('disabled', readonly);
                    row.find('.btn-school-delete').click((function(schoolId, listIndex) {
                        return function() {
                            alert('确定删除吗？');
                            var countryListIndex = $('#school-page-country-select').val();
                            var data = {
                                countryListIndex: countryListIndex - 1,
                                schoolId: schoolId,
                                schoolListIndex: listIndex - 1,
                            }
                            $.ajax({
                                url: '/school/delete',
                                type: 'DELETE',
                                headers: {
                                    'Content-Type': 'application/json'
                                },
                                data: JSON.stringify(data),
                                success: function(data) {
                                    if (table.children().length === 1) {
                                        $('#school-pagination').children().last().remove();
                                        if (currentSchoolPage > 1) {
                                            currentSchoolPage--;
                                        }
                                    }
                                    fetchSchoolData(currentSchoolPage);
                                }
                            });
                        }
                    })(school[i].schoolId, listIndex));

                    var countryListIndex = $('#school-page-country-select').val();
                    row.find('.btn-item').click((function(schoolListIndex, countryListIndex) {
                        return function() {
                            $("#manage-school-content").css("display", "none");
                            $("#manage-item-content").css("display", "block");
                            $("#manage-item-content").css({
                                "position": "absolute", // 使用绝对定位
                                "top": "80px", 
                                "left": "20vw" 
                            });
                            initItem(countryListIndex, schoolListIndex);
                        }
                    })(listIndex, countryListIndex));
                }

                $('#school-pagination').empty();
                for (let i = 1; i <= totalSchoolPage; i++) {
                    $('#school-pagination').append(`<a class="page-link" href="#" data-page="${i}">${i}</a>`);
                }
                $('#add-school-btn').prop('disabled', readonly);
            }
        });
        
    }

    $(document).on('click', '#school-pagination .page-link', function(e) {
        e.preventDefault();
        currentSchoolPage = $(this).data('page');
        listIndex = $('#school-page-country-select').val();
        fetchSchoolData(listIndex, currentSchoolPage);
    });

    $("#add-school-btn").click(function() {
        var listIndex = $('#school-page-country-select').val();
        if (currentSchoolPage != totalSchoolPage) {
            currentSchoolPage = totalSchoolPage;
            fetchSchoolData(listIndex, currentSchoolPage);
        }
        $.ajax({
            url: 'school/create',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify({
                countryListIndex: listIndex - 1,
                pageNum: pageNum,
            }),
            success: function(data) {
                fetchSchoolData(listIndex, data.totalPage);
            }
        });
    });

    function schoolTextChange(schoolId, listIndex, field, value) {
        var countryListIndex = $('#school-page-country-select').val();
        $.ajax({
            url: '/school/change',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify({
                countryListIndex: countryListIndex - 1,
                schoolId: schoolId,
                schoolListIndex: listIndex - 1,
                updateField: field,
                updateValue: value
            }),
            success: function(data) {
            }
        });
    }

    function initItem(countrylistIndex = 0, schoolListIndex = 0){
        $.ajax({
            url: '/school/initPage',
            type: 'GET',
            success: function(data) {
                $('#item-table tbody').empty();
                $('#add-item-btn').prop('disabled', true);
                var itemSwitch = $('#item-switch input[type="checkbox"]');
                itemSwitch.prop('checked', false);
                itemSwitch.change(function() {
                    var readonly = !$(this).prop('checked');
                    $('#add-item-btn').prop('disabled', readonly);
                    $('.input-text').prop('readonly', readonly);
                    $('.input-select').prop('disabled', readonly);
                });
                var allCountry = data.results;
                var countrySelect = $('#item-page-country-select');
                countrySelect.empty();
                countrySelect.append('<option value="0">请选择国家</option>');
                for (var i = 0; i < allCountry.length; i++) {
                    var option = $(`<option value="${i+1}">${allCountry[i]}</option>`);
                    countrySelect.append(option);
                }
                countrySelect.off('change').change(function() {
                    fetchSchoolList(schoolListIndex);
                });
                if(countrylistIndex != 0){
                    countrySelect.val(countrylistIndex).trigger('change');
                }
            }
        });
    }

    function fetchSchoolList(schoolListIndex = 0){
        var countryListIndex = $('#item-page-country-select').val();
        $.get("/item/getSchool", {countryListIndex: countryListIndex-1}, function(data) {
            var allSchool = data.results;
            var schoolSelect = $('#item-page-school-select');
            schoolSelect.empty();
            schoolSelect.append('<option value="0">请选择学校</option>');
            for (var i = 0; i < allSchool.length; i++) {
                var option = $(`<option value="${i+1}">${allSchool[i]}</option>`);
                schoolSelect.append(option);
            }
            schoolSelect.off('change').change(function() {
                var schoolListIndex = $(this).val();
                fetchItemData(schoolListIndex);
            });
            if(schoolListIndex != 0){
                schoolSelect.val(schoolListIndex).trigger('change');
            }
        });
    }

    function fetchItemData(schoolListIndex, page = 1){
        var countryListIndex = $('#item-page-country-select').val();
        var data = {
            schoolListIndex : schoolListIndex - 1,
            countryListIndex : countryListIndex - 1,
            page: page,
            pageNum: pageNum
        }
        $.ajax({
            url: '/school/editItem',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify(data),
            success: function(data) {
                var item = data.results;
                totalItemPage = data.totalPage;
                var table = $('#item-table tbody');
                table.empty();
                var itemSwitch = $('#item-switch input[type="checkbox"]');
                var readonly = !itemSwitch.prop('checked');

                for(let i=0; i<item.length; i++){
                    var listIndex = (page - 1) * pageNum + i + 1;
                    var itemName = $(`<input type="text" class="input-text" value="${item[i].itemName}" />`);
                    var levelDescription = $(`<input type="text" class="input-text" value="${item[i].levelDescription}" />`);
                    var itemRemark = $(`<input type="text" class="input-text" value="${item[i].itemRemark}" />`);
                    var row = $(
                        `<tr>
                            <td>${listIndex}</td>
                            <td>${itemName.prop('outerHTML')}</td>
                            <td>${levelDescription.prop('outerHTML')}</td>
                            <td>${itemRemark.prop('outerHTML')}</td>
                            <td>${item[i].levelRate.length}</td>
                            <td>
                                <a href=# class="btn btn-item-edit">编辑</a>
                                <a href=# class="btn btn-item-delete">删除</a>
                            </td>
                        </tr>`
                    );
                    table.append(row);
                    row.find('input.input-text').eq(0).change((function(itemId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            itemTextChange(itemId, listIndex, field, value);
                        }
                    })(item[i].itemId, listIndex, 'itemName'));

                    row.find('input.input-text').eq(1).change((function(itemId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            itemTextChange(itemId, listIndex, field, value);
                        }
                    })(item[i].itemId, listIndex, 'levelDescription'));

                    row.find('input.input-text').eq(2).change((function(itemId, listIndex, field) {
                        return function() {
                            var value = $(this).val();
                            itemTextChange(itemId, listIndex, field, value);
                        }
                    })(item[i].itemId, listIndex, 'itemRemark'));

                    row.find('input.input-text').prop('readonly', readonly);

                    row.find('.btn-item-edit').off('click').click((function(countryListIndex, schoolListIndex, itemListIndex, currentItemPage) {
                        return function() {
                            disableSidebarButton(); // 禁用侧边栏按钮
                            var data = {
                                countryListIndex: countryListIndex - 1,
                                schoolListIndex: schoolListIndex - 1,
                                itemListIndex: itemListIndex - 1
                            }
                            $.ajax({
                                url: '/item/changeProportion/show',
                                type: 'POST',
                                headers: {
                                    'Content-Type': 'application/json'
                                },
                                data: JSON.stringify(data),
                                success: function(response) {
                                    $('#item-model').css('display', 'block');
                                    $("#manage-item-content").css('pointer-events', 'none');
                                    $("#manage-item-content").css('opacity', '0.5');
                                    $('#item-name-input').val(response.item.itemName);
                                    $('#level-description-input').val(response.item.levelDescription);
                                    $('#item-remark-input').val(response.item.itemRemark);
                                    $('#save-item-btn').off('click').click((function(countryListIndex, schoolListIndex, itemListIndex, itemId){
                                        return function(){
                                            enableSidebarButton(); // 启用侧边栏按钮
                                            var levelRate = getLevelRateData();
                                            var itemName = $('#item-name-input').val();
                                            var levelDescription = $('#level-description-input').val();
                                            var itemRemark = $('#item-remark-input').val();
                                            var data = {
                                                countryListIndex: countryListIndex-1,
                                                schoolListIndex: schoolListIndex-1,
                                                itemListIndex: itemListIndex-1,
                                                itemId: itemId,
                                                itemName: itemName,
                                                levelDescription: levelDescription,
                                                itemRemark: itemRemark,
                                                levelRate: levelRate
                                            }
                                            $.ajax({
                                                url: '/item/changeProportion/save',
                                                type: 'POST',
                                                headers: {
                                                    'Content-Type': 'application/json'
                                                },
                                                data: JSON.stringify(data),
                                                success: function(data) {
                                                    $('#item-model').css('display', 'none');
                                                    $("#manage-item-content").css('pointer-events', 'auto');
                                                    $("#manage-item-content").css('opacity', '1');
                                                    fetchItemData(schoolListIndex, currentItemPage);
                                                }
                                            });
                                        }
                                    })(countryListIndex, schoolListIndex, itemListIndex, response.item.itemId));
                                    fetchItemWindowData(response.item.levelRate);
                                }
                            });
                        }
                    })(countryListIndex, schoolListIndex, listIndex, currentItemPage));

                    row.find('.btn-item-delete').click((function(itemId, listIndex) {
                        return function() {
                            alert('确定删除吗？');
                            var countryListIndex = $('#item-page-country-select').val();
                            var schoolListIndex = $('#item-page-school-select').val();
                            var data = {
                                countryListIndex: countryListIndex - 1,
                                schoolListIndex: schoolListIndex - 1,
                                itemId: itemId,
                                listIndex: listIndex - 1
                            }
                            $.ajax({
                                url: '/item/delete',
                                type: 'DELETE',
                                headers: {
                                    'Content-Type': 'application/json'
                                },
                                data: JSON.stringify(data),
                                success: function(data) {
                                    if (table.children().length === 1) {
                                        $('#item-pagination').children().last().remove();
                                        if (currentItemPage > 1) {
                                            currentItemPage--;
                                        }
                                    }
                                    fetchItemData(schoolListIndex, currentItemPage);
                                }
                            });
                        }
                    })(item[i].itemId, listIndex));
                }

                $('#item-pagination').empty();
                for (let i = 1; i <= totalItemPage; i++) {
                    $('#item-pagination').append(`<a class="page-link" href="#" data-page="${i}">${i}</a>`);
                }
                $('#add-item-btn').prop('disabled', readonly);
            }
        });
    }

    $(document).on('click', '#item-pagination .page-link', function(e) {
        e.preventDefault();
        currentItemPage = $(this).data('page');
        var schoolListIndex = $('#item-page-school-select').val();
        fetchItemData(schoolListIndex, currentItemPage);
    });

    $("#add-item-btn").click(function() {
        var countryListIndex = $('#item-page-country-select').val();
        var schoolListIndex = $('#item-page-school-select').val();
        if (currentItemPage != totalItemPage) {
            currentItemPage = totalItemPage;
            fetchItemData(schoolListIndex, currentItemPage);
        }
        $.ajax({
            url: 'item/create',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify({
                countryListIndex: countryListIndex - 1,
                schoolListIndex: schoolListIndex - 1,
                pageNum: pageNum,
            }),
            success: function(data) {
                fetchItemData(schoolListIndex, data.totalPage);
            }
        });
    });

    function itemTextChange(itemId, listIndex, field, value) {
        var countryListIndex = $('#item-page-country-select').val();
        var schoolListIndex = $('#item-page-school-select').val();
        $.ajax({
            url: "/item/change",
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify({
                countryListIndex: countryListIndex-1,
                schoolListIndex: schoolListIndex-1,
                itemId: itemId,
                listIndex: listIndex - 1,
                updateField: field,
                updateValue: value
            }),
            success: function(data){
            }
        });
    }

    function fetchItemWindowData(levelRate){
        var table = $('#level-table tbody');
        table.empty();
        for (let i=0; i<levelRate.length; i++) {
            var levelId = $(`<input type="text" class="input-text" value="${levelRate[i].levelId}" />`);
            var rate = $(`<input type="text" class="input-text" value="${levelRate[i].levelRate}" />`);
            var checkbox = $(`<input type="checkbox" id="check-box-${i}" />`).attr('checked', levelRate[i].ifNotCombine);
            var row = $(
                `<tr>
                    <td>${levelId.prop('outerHTML')}</td>
                    <td>${rate.prop('outerHTML')}</td>
                    <td>${checkbox.prop('outerHTML')}</td>
                    <td>
                        <a href=# class="btn btn-level-delete">删除</a>
                    </td>
                </tr>`
            );
            row.find('.btn-level-delete').off('click').click((function(i){
                return function() {
                    var levelRate = getLevelRateData();
                    levelRate.splice(i, 1);
                    fetchItemWindowData(levelRate);
                }
            })(i));
            table.append(row);
        }
    }

    function getLevelRateData(){
        var levelRate = [];
        $('#level-table tbody tr').each(function() {
            var levelId = parseInt($(this).find('input.input-text').eq(0).val());
            var rate = $(this).find('input.input-text').eq(1).val();
            var ifNotCombine = $(this).find('input[type="checkbox"]').prop('checked');
            levelRate.push({levelId: levelId, levelRate: rate, ifNotCombine: ifNotCombine});
        });
        return levelRate;
    }

    $('#add-level-btn').click(function(){
        var levelRate = getLevelRateData();
        levelRate.push({levelId: 0, levelRate: 0});
        fetchItemWindowData(levelRate);
    })

    $('#cancel-item-btn').click(function(){
        enableSidebarButton(); // 启用侧边栏按钮
        $('#item-model').css('display', 'none');
        $("#manage-item-content").css('pointer-events', 'auto');
        $("#manage-item-content").css('opacity', '1');
    });

    function initUser(){
        $('#add-user-btn').prop('disabled', true);
        var userSwitch = $('#user-switch input[type="checkbox"]');
        userSwitch.prop('checked', false);
        userSwitch.change(function() {
            var readonly = !$(this).prop('checked');
            $('#add-user-btn').prop('disabled', readonly);
            $('.input-text').prop('readonly', readonly);
        });
        fetchUserData();
    }

    function fetchUserData(page = currentUserPage){
        var data = {
            page: page,
            pageNum: pageNum
        }
        $.ajax({
            url: '/user/show',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify(data),
            success: function(data) {
                totalUserPage = data.totalPage;
                var table = $('#user-table tbody');
                table.empty();
                var userSwitch = $('#user-switch input[type="checkbox"]');
                var readonly = !userSwitch.prop('checked');

                for(let i = 0; i < data.results.length; ++i){
                    var listIndex = (page - 1) * pageNum + i + 1;
                    var userAccount = $(`<input type="text" class="input-text" value="${data.results[i].userAccount}" readonly />`);
                    var userEmail = $(`<input type="text" class="input-text" value="${data.results[i].userEmail}" readonly />`);
                    var userNumber = $(`<input type="text" class="input-text" value="${data.results[i].userNumber}" readonly />`);
                    var userLevel = $(`<input type="text" class="input-text" value="${data.results[i].userLevel}" readonly />`);
                    var studentCount = $(`<input type="text" class="input-text" value="${data.results[i].studentCount}" readonly />`);
                    var userPasswd = $(`<input type="text" class="input-text" value="${data.results[i].userPasswd}" readonly />`);
                    var row = $(
                        `<tr>
                            <td>${listIndex}</td>
                            <td>${userAccount.prop('outerHTML')}</td>
                            <td>${userEmail.prop('outerHTML')}</td>
                            <td>${userNumber.prop('outerHTML')}</td>
                            <td>${userLevel.prop('outerHTML')}</td>
                            <td>${studentCount.prop('outerHTML')}</td>
                            <td>${userPasswd.prop('outerHTML')}</td>
                            <td>
                                <a href=# class="btn btn-delete">删除</a>
                            </td>
                        </tr>`
                    );
                    table.append(row);
                    row.find('input.input-text').eq(0).change((function(userId, listIndex, inputText) {
                        return function() {
                            var value = $(this).val();
                            userTextChange(userId, listIndex, inputText, value);
                        }
                    })(data.results[i].userId, listIndex, 'userAccount'));
                    row.find('input.input-text').eq(1).change((function(userId, listIndex, inputText) {
                        return function() {
                            var value = $(this).val();
                            userTextChange(userId, listIndex, inputText, value);
                        }
                    })(data.results[i].userId, listIndex, 'userEmail'));
                    row.find('input.input-text').eq(2).change((function(userId, listIndex, inputText) {
                        return function() {
                            var value = $(this).val();
                            userTextChange(userId, listIndex, inputText, value);
                        }
                    })(data.results[i].userId, listIndex, 'userNumber'));
                    row.find('input.input-text').eq(3).change((function(userId, listIndex, inputText) {
                        return function() {
                            var value = parseInt($(this).val());
                            userTextChange(userId, listIndex, inputText, value);
                        }
                    })(data.results[i].userId, listIndex, 'userLevel'));
                    row.find('input.input-text').eq(4).change((function(userId, listIndex, inputText) {
                        return function() {
                            var value = parseInt($(this).val());
                            userTextChange(userId, listIndex, inputText, value);
                        }
                    })(data.results[i].userId, listIndex, 'studentCount'));
                    row.find('input.input-text').eq(5).change((function(userId, listIndex, inputText) {
                        return function() {
                            var value = $(this).val();
                            userTextChange(userId, listIndex, inputText, value);
                        }
                    })(data.results[i].userId, listIndex, 'userPasswd'));
                    row.find('input.input-text').prop('readonly', readonly);
                    row.find('.btn-delete').click((function(userId, listIndex) {
                        return function() {
                            alert('确定删除吗？')
                            var data = {
                                userId: userId,
                                listIndex: listIndex - 1,
                            }
                            $.ajax({
                                url: '/user/delete',
                                type: 'DELETE',
                                headers: {
                                    'Content-Type': 'application/json'
                                },
                                data: JSON.stringify(data),
                                success: function(data) {
                                    if (table.children().length === 1) {
                                        $('#user-pagination').children().last().remove();
                                        if (currentUserPage > 1) {
                                            currentUserPage--;
                                        }
                                    }
                                    fetchUserData(currentUserPage);
                                }
                            });
                        }
                    })(data.results[i].userId, listIndex));
                }
                $('#user-pagination').empty();
                for (let i = 1; i <= data.totalPage; i++) {
                    $('#user-pagination').append(`<a class="page-link" href="#" data-page="${i}">${i}</a>`);
                }
                $('#add-user-btn').prop('disabled', readonly);
            }
        })
    }

    $(document).on('click', '#user-pagination .page-link', function(e) {
        e.preventDefault();
        currentUserPage = $(this).data('page');
        fetchUserData(currentUserPage);
    });

    $("#add-user-btn").click(function() {
        if (currentUserPage != totalUserPage) {
            currentUserPage = totalUserPage;
            fetchCountryData(currentUserPage);
        }
        $.get("/user/create", {pageNum: pageNum}, function(response) {
            currentUserPage = response.totalPage;
            fetchUserData(response.totalPage);
        });
    });

    function userTextChange(userId, listIndex, field, value){
        $.ajax({
            url: '/user/changeUser',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify({
                userId: userId,
                listIndex: listIndex - 1,
                updateField: field,
                updateValue: value
            }),
            success: function(data) {
            }
        })

    }

    function initSystemSet(){
        $('#add-type-btn').prop('disabled', true);
        var typeSwitch = $('#system-switch input[type="checkbox"]');
        typeSwitch.prop('checked', false);
        typeSwitch.change(function() {
            var readonly = !$(this).prop('checked');
            $('#add-type-btn').prop('disabled', readonly);
            $('.input-text').prop('readonly', readonly);
        });
        fetchSystemData();
    }

    function fetchSystemData() {
        $.get("/system/show", function(response) {
            var system = response.system;
            var maxUserLevelText = $('#user-level-input');
            maxUserLevelText.val(system.maxUserLevel);
            maxUserLevelText.change((function(listIndex, field) {
                return function() {
                    var value = $(this).val();
                    systemTextChange(listIndex, field, parseInt(value));
                }
            })(-1, 'maxUserLevel'));

            var table = $('#system-table tbody');
            table.empty();
            var systemSwitch = $('#system-switch  input[type="checkbox"]');
            var readonly = !systemSwitch.prop('checked');

            for (let i = 0; i < system.schoolTypeList.length; i++) {
                var listIndex = i + 1;
                var schoolTypeNameText =  $(`<input type="text" class="input-text" value="${system.schoolTypeList[i].schoolTypeName}" readonly />`);
                var row = $(
                    `<tr>
                        <td>${listIndex}</td>
                        <td>${schoolTypeNameText.prop('outerHTML')}</td>
                        <td>
                            <a href=# class="btn btn-delete">删除</a>
                        </td>
                    </tr>`
                );
                table.append(row);
                row.find('input.input-text').eq(0).change((function(listIndex, field) {
                    return function() {
                        var value = $(this).val();
                        systemTextChange(listIndex, field, value);
                    }
                })(listIndex, 'schoolTypeList'));
                row.find('input.input-text').prop('readonly', readonly);
                row.find('.btn-delete').off('click').click((function(listIndex) {
                    return function() {
                        alert('确定删除吗？');
                        $.ajax({
                            url: "/system/delete",
                            type: "DELETE",
                            headers: {
                                'Content-Type': 'application/json'
                            },
                            data: JSON.stringify({listIndex: listIndex - 1}),
                            statusCode: {
                                200: function(response) {
                                    fetchSystemData();
                                },
                                400: function(jqXHR) {
                                    var response = JSON.parse(jqXHR.responseText);
                                    alert(response.usedSchool);
                                }
                            }
                        });
                    }
                })(listIndex));
            }
        });
    }

    $("#add-type-btn").click(function() {
        $.get("/system/create", {pageNum: pageNum}, function(response) {
            fetchSystemData();
        });
    });

    function systemTextChange(listIndex, field, value) {
        $.ajax({
            url: '/system/change',
            type: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify({
                listIndex: listIndex - 1,
                updateField: field,
                updateValue: value
            }),
            success: function(response) {

            }
        });
    }

    function disableSidebarButton(){
        $(".button-list-button").prop("disabled", true);
    }

    function enableSidebarButton(){
        $(".button-list-button").prop("disabled", false);
    }

    fetchCountryData();
    initCountry();
});