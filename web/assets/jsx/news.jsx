/**
 * Created by myan on 21/09/15.
 */

var apiKeyPostfix = "key=abc123";
var editable = false;
var apiHost = "http://writeuptube.com:5050";

var getURLParameter = function (sParam) {
  var sPageURL = window.location.search.substring(1);
  var sURLVariables = sPageURL.split('&');
  for (var i = 0; i < sURLVariables.length; i++) {
    var sParameterName = sURLVariables[i].split('=');
    if (sParameterName[0] == sParam) {
      return sParameterName[1];
    }
  }
};

var prevPages = function () {
  $("ul.pagination>li.pageno>a").map(function () {
    pageno = parseInt($(this).text());
    if (pageno > 5) {
      $(this).text(pageno - 5);
    }
  });
};

var nextPages = function () {
  $("ul.pagination>li.pageno>a").map(function () {
    pageno = parseInt($(this).text());
    $(this).text(pageno + 5);
  });
};


var ArticlePane = React.createClass({
  getInitialState: function () {
    return {data: []};
  },

  componentDidMount: function () {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function (data) {
        this.setState({data: data})
      }.bind(this),
      error: function (xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },

  goActivePage: function () {
    var ap = parseInt($(window.event.srcElement).text());
    var path = "news.html?page=" + ap;
    var editable = getURLParameter("edit")
    if (editable == "1" || editable == "true") {
      path += "&edit=1";
    }
    var url = apiHost + "/articles?page=" + (ap - 1) + "&" + apiKeyPostfix;

    $("ul.pagination>li.pageno").removeClass("active");
    history.pushState({name: "readcn"}, "ReadCN", path);

    $.ajax({
      url: url,
      dataType: 'json',
      cache: false,
      success: function (data) {
        this.setState({data: data})
      }.bind(this),
      error: function (xhr, status, err) {
        console.error(url, status, err.toString());
      }.bind(this)
    });
  },

  render: function () {
    return (
      <div className="articles col-md-9 col-sm-9 blog-posts">
        <ArticleList data={this.state.data}/>
        <ul className="pagination">
          <li><a href="javascript: prevPages();">Prev</a></li>
          <li className="pageno active"><a href="javascript:;" onClick={this.goActivePage}>1</a></li>
          <li className="pageno"><a href="javascript:;" onClick={this.goActivePage}>2</a></li>
          <li className="pageno"><a href="javascript:;" onClick={this.goActivePage}>3</a></li>
          <li className="pageno"><a href="javascript:;" onClick={this.goActivePage}>4</a></li>
          <li className="pageno"><a href="javascript:;" onClick={this.goActivePage}>5</a></li>
          <li><a href="javascript: nextPages();">Next</a></li>
        </ul>

      </div>
    )
  }
});


var ArticleList = React.createClass({
  render: function () {
    var articleNodes = this.props.data.map(function (article) {
      return (
        <ArticleSummary data={article}/>
      );
    });
    return (
      <div className="article-list row">
        {articleNodes}
      </div>
    );
  }
});

var ArticleSummary = React.createClass({
  render: function () {
    var rawMarkup = this.props.data.preview;
    var genlink = function (permalink) {
      return '/static/read.html?v=' + permalink;
    };
    var editlink = function (permalink) {
      return '/static/edit.html?v=' + permalink;
    };
    var insertEditLink = function () {
      if (editable) {
        return (
          <span>
            <a href="#" className="more"> | </a><a href={editlink(this.props.data.permalink)} className="more"> Edit <i
            className="icon-angle-right"></i></a>
          </span>
        );
      }
    }.bind(this);

    return (
      <div>
        <h4><a href={genlink(this.props.data.permalink)}>{this.props.data.title}</a></h4>
        <ul className="blog-info">
          <li><i className="fa fa-calendar"></i>{this.props.data.publishedAt}</li>
          <li><i className="fa fa-comments"></i>{this.props.data.num_comment}</li>
          <li><i className="fa fa-tags"></i>{this.props.data.tags}</li>
        </ul>
        <pre>
          <div dangerouslySetInnerHTML={{__html: rawMarkup}}/>
        </pre>
        <a href={genlink(this.props.data.permalink)} className="more">Read more <i className="icon-angle-right"></i></a>
        {insertEditLink()}
        <hr className="blog-post-sep"/>
      </div>
    );
  }
});

$(function () {
  var edit = getURLParameter("edit");
  if (edit == "true" || edit == "1") {
    editable = true;
  }


  activePageNumber = parseInt(getURLParameter("page"));
  if (isNaN(activePageNumber)) {
    activePageNumber = 1;
  }

  var apiStartURL = function () {
    return apiHost + "/articles?" + "page=" + (activePageNumber - 1) + "&" + apiKeyPostfix;
  };

  React.render(
    //<ArticlePane data={data} />,
    <ArticlePane url={apiStartURL()}/>,
    document.getElementById('articles')
  );

  startPageNo = Math.ceil(activePageNumber / 5) * 5 - 4;
  $("ul.pagination>li.pageno>a").map(function () {
    $(this).text(startPageNo++);
  });
  $("ul.pagination>li.pageno").removeClass("active");
  activePageLable = $("ul.pagination>li.pageno").get((activePageNumber - 1) % 5);
  $(activePageLable).addClass("active");
});


