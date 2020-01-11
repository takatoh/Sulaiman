let app;
let photoList;
let triggerFlag = false;

$(function() {
  let upload_dialog;
  let delete_dialog;

  function upload(event) {
    let ext = $("#upload-form input[name=file]").val().split(".").pop().toLowerCase();
    if ($.inArray(ext, ["jpg", "jpeg", "png", "gif"]) == -1) {
      alert("Unsupported format!\nSupport only jpg, png, gif.");
      clear_upload_vals();
    } else {
      let fd = new FormData($("#upload-form").get(0));
      $.ajax({
        url: "/upload",
        type: "POST",
        data: fd,
        contentType: false,
        processData: false,
        dataType: "json"
      }).done(function(response) {
        clear_upload_vals();
        app.photoList.unshift(response.photo);
        if (response.delete_photo_id > 0) {
          app.photoList = app.photoList.filter(function(p){
            return p.id != response.delete_photo_id;
          });
        }
        upload_dialog.dialog("close");
      });
    }
  }

  function clear_upload_vals() {
    $("input[name=file]").val("");
    $("input[name=key]").val("");
  }

  function delete_photo(event) {
    let fd = new FormData($("#delete-form").get(0));
    $.ajax({
      url: "/delete",
      type: "DELETE",
      data: fd,
      contentType: false,
      processData: false,
      dataType: "json"
    }).done(function(response) {
      if (response.status == "OK") {
        clear_delete_vals();
        photoList.some(function(v, i) {
          if (v.id == response.photo_id) {
            photoList.splice(i, 1);
          }
        });
        alert("Deleted: " + response.photo_id);
      } else {
        clear_delete_vals();
        alert("Error! CAN'T delete: " + response.photo_id);
      }
      delete_dialog.dialog("close");
    });
  }

  function clear_delete_vals() {
    $("#delete-form input[name=id]").val("");
    $("#delete-form input[name=key]").val("");
  }

  $.ajax({
    type: "GET",
    url: "/title",
    dataType: "text"
  }).done(function(response) {
    $("head title").text(response)
    $("h1 a").text(response);
  });

  let next_url = $("#next-link").attr("href");
  $.ajax({
    type: "GET",
    url: next_url,
    dataType: "json"
  }).done(function(response) {
    if (response.photos) {
      photoList = response.photos;
    } else {
      photoList = [];
    }
    app = new Vue({
      el: "#content",
      data: {
        photoList: photoList
      }
    });
    if (response.next) {
      $("#next-link").attr("href", response.next);
    } else {
      $("#next-link").remove();
    }
  });

  upload_dialog = $("#upload-dialog").dialog({
    autoOpen: false,
    modal: true,
    draggable: false,
    resizable: false,
    width: 350,
    buttons: {
      Upload: upload,
      Cancel: function() {
        clear_upload_vals();
        upload_dialog.dialog("close");
      }
    }
  });

  $("#upload").button().on("click", function() {
    upload_dialog.dialog("open");
  });

  delete_dialog = $("#delete-dialog").dialog({
    autoOpen: false,
    modal: true,
    draggable: false,
    resizable: false,
    width: 300,
    buttons: {
      Delete: delete_photo,
      Cancel: function() {
        clear_delete_vals();
        delete_dialog.dialog("close");
      }
    }
  });

  $("#delete").button().on("click", function() {
    delete_dialog.dialog("open");
  });

  $(window).on("load scroll", function() {
    let documentHeight = $(document).height();
    let scrollBottomPosition = $(window).height() + $(window).scrollTop();
    let triggerPoint = documentHeight - scrollBottomPosition;
    if (!triggerFlag && triggerPoint <= 50) {
      triggerFlag = true;
      let next_url = $("#next-link").attr("href");
      $.ajax({
        type: "GET",
        url: next_url,
        dataType: "json"
      }).done(function(response) {
        if (response.photos) {
          response.photos.forEach(function(v) { app.photoList.push(v); });
        }
        if (response.next) {
          $("#next-link").attr("href", response.next);
        } else {
          $("#next-link").remove();
        }
      });
    }
    if (triggerPoint > 50) {
      triggerFlag = false;
    }
  });
});
