var ulFather = document.getElementById("ulFather");
var box = document.getElementById("box");
let current_page = 1;
var setTotalCount = 500;
var maxPage = 30;


function restore(){
    document.body.removeChild(document.getElementById("overlay"));
    document.body.removeChild(document.getElementById("expand"));
}

function expandPhoto(){
    var overlay = document.createElement("div");
    overlay.setAttribute("id","overlay");
    overlay.setAttribute("class","overlay");
    document.body.appendChild(overlay);

    var img = document.createElement("img");
    img.setAttribute("id","expand")
    img.setAttribute("class","overlayimg");
    img.src = this.getAttribute("src");
    document.getElementById("overlay").appendChild(img);

    img.onclick = restore;
}


function addExpand() {
    var imgs = document.getElementsByTagName("img");
    imgs[0].focus();
    for(var i = 0;i<imgs.length;i++){
        imgs[i].onclick = expandPhoto;
        imgs[i].onkeydown = expandPhoto;
    }
}

//将input时间区间最晚时间默认值设置为现在
$(document).ready(function () {
    var time = new Date();
    var day = ("0" + time.getDate()).slice(-2);
    var month = ("0" + (time.getMonth() + 1)).slice(-2);
    var today = time.getFullYear() + "-" + (month) + "-" + (day);
    $('#timeOver').val(today);
})


// 请求到数据 展示
function displayResponseData(arr, isInit) {
    setTotalCount = arr.total;
    maxPage = Math.ceil(arr.total / 20);

    $('#box').paging({
        initPageNo: isInit ? 1 : current_page, // 初始页码
        totalPages: maxPage, //总页数
        totalCount: '合计' + setTotalCount + '条数据', // 条目总数
        slideSpeed: 600, // 缓动速度。单位毫秒
        jump: true, //是否支持跳转
        callback: function (page) { // 回调函数
            console.log(page);
            current_page = page;
        }
    })
    $("#search,#prePage,#nextPage,#lastPage,#jumpBtn,#pageSelect").click(function () {
        getRecords(false);
    })
    //console.log(setTotalCount);
    var childs = ulFather.childNodes;
    for (var i = childs.length - 1; i >= 0; i--) {
        ulFather.removeChild(childs[i]);
    }
    for (const pic of arr.data.pics) {
        var liC = document.createElement("li");
        var Img = document.createElement("img");
        var breakType = document.createElement("div");
        var location = document.createElement("div");
        var carNum = document.createElement("div");
        var breakTime = document.createElement("div");
        Img.src = pic.pic_path;
        console.log(pic)
        Img.className = 'showImg';
        breakType.innerText = "违法类型: " + pic.rule_type;
        location.innerText = "违法地点: " + pic.location;
        carNum.innerText = "车牌号码: " + pic.lic_plate;
        breakTime.innerText = "拍摄时间: " + pic.shoot_time;
        ulFather.appendChild(liC);
        liC.appendChild(Img);
        liC.appendChild(location);
        liC.appendChild(breakType);
        liC.appendChild(carNum);
        liC.appendChild(breakTime);
    }
    addExpand();
}

function handleSubmit() {
    // 取消默认的提交表单触发事件
    return false;
}

function getRecords(isInit) {
    var formObject = {};
    var formArray = $("#form1").serializeArray();
    var oTimer = formArray[1];
    formArray[1].value = new Date(oTimer.value).getTime() / 1000;
    oTimer = formArray[2];
    formArray[2].value = new Date(oTimer.value).getTime() / 1000;
    formArray[3] = { name: 'page', value: current_page };
    formArray[4] = { name: 'page_size', value: 20 };
    console.log(current_page)
    $.each(formArray, function (i, item) {
        formObject[item.name] = item.value;
    });
    console.log(formObject);
    $.ajax({
        url: "http://localhost:8081/api/v1/records",
        type: "post",
        hearders: { "contentType": "application/json; charset=utf-8" },
        data: JSON.stringify(formObject),
        //data:formArray,
        dataType: "json",
        success: function (data) {
            //var cData = JSON.parse(data);
            displayResponseData(data, isInit);
            console.log(data);
        },
        error: function (e) {
            console.log("错误");
        }
    });
}

$("#search").on("click", function () {
    getRecords(true);
});

$("body").ready(function (){
    getRecords(true)
})