/**
 * Created by myan on 21/09/15.
 */


var editable = false;

var searchTarget = function () {
  var char = getURLParameter("c");
  var vocabulary = getURLParameter("w");
  var grammar = getURLParameter("m");
  var target = char;
  if (typeof(vocabulary) != "undefined") {
    target = vocabulary;
  } else if (typeof(grammar) != "undefined") {
    target = grammar;
  }
  return decodeURI(target);
};


var SearchView = React.createClass({
  getInitialState: function () {
    return {
      data: [],
      activeLevel: 1,
      startPage: 1,
      numPages: 5,
      currentPage: 1
    };
  },

  componentDidMount: function () {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({data: data})
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },

  componentDidUpdate: function () {
    $("div.search-result").map(function () {
      var hl = searchTarget();
      var html = $(this).html();
      var re = new RegExp(hl, 'g');
      var newhtml = html.replace(re, '<span class="highlighted">' + hl + '</span>');
      $(this).html(newhtml);
    });
  },

  switchLevel: function () {
    var g = $(window.event.srcElement).text().substring(6);
    var char = getURLParameter("c");
    var vocabulary = getURLParameter("w");
    var grammar = getURLParameter("m");

    var path = "/bychar?t=" + char;
    if (typeof(vocabulary) != "undefined") {
      path = "/byvocabulary?t=" + vocabulary;
    } else if (typeof(grammar) != "undefined") {
      path = "/bygrammar?t=" + grammar;
    }
    var url = apiHost + path + "&p=0&g=" + g + "&" + apiKeyPostfix;

    $.ajax({
      url: url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({
          activeLevel: g,
          currentPage: 1,
          data: data
        })
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(url, status, err.toString());
      }.bind(this)
    });

  },

  setPages: function (startPage, numPages, currentPage) {
    this.setState({
      startPage: startPage,
      numPages: numPages,
      currentPage: currentPage
    });
  },

  searchTarget: function (evt) {
    var target = $("input#search-target").val();
    var q = '?c=';
    if (target.length > 1) {        
        if (target.split('/').length >= 2) {          
          q = '?m=';
        } else {
          q = '?w='; 
        }        
    }
    var url = window.location.href;
    var newURL = url.replace(/(\?[cwm])\=([^\&]+)/, q + encodeURIComponent(target));
    window.location.href = newURL;
    evt.preventDefault();
  },

  gotoPage: function () {
    var pg = parseInt($(window.event.srcElement).text(), 10);
    var char = getURLParameter("c");
    var vocabulary = getURLParameter("w");
    var grammar = getURLParameter("m");
    var g = this.state.activeLevel;

    var path = "/bychar?t=" + char;
    if (typeof(vocabulary) != "undefined") {
      path = "/byvocabulary?t=" + vocabulary;
    } else if (typeof(grammar) != "undefined") {
      path = "/bygrammar?t=" + grammar;
    }

    var url = apiHost + path + "&p=" + (pg-1) + "&g=" + g + "&" + apiKeyPostfix;
    $.ajax({
      url: url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({
          currentPage: pg,
          data: data
        })
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(url, status, err.toString());
      }.bind(this)
    });

  },

  render: function () {
    return (
      <div className="col-md-12">
        <div className="content-page">
          <form className="content-search-view2" action="#">
            <div className="input-group">
              <input id="search-target" type="text" className="form-control" placeholder="Search..." />
              <span className="input-group-btn">
                <button className="btn btn-primary" onClick={this.searchTarget}>Search</button>
              </span>
            </div>
          </form>
        </div>

        <LevelSwitchTab levels={this.props.num_levels} activeLevel={this.state.activeLevel}
                        levelLabel={this.state.levelLabel} onSwitch={this.switchLevel} />

        <hr/>
        <SearchResultList data={this.state.data} />

        <PaginationTab startPage={this.state.startPage} numPages={this.state.numPages}
                       currentPage={this.state.currentPage} setPages={this.setPages}
                       gotoPage={this.gotoPage} />
      </div>
    );
  }
});

var SearchResultList = React.createClass({
  render: function () {
    var searchResults = this.props.data.map(function (sr) {
      return (<SearchResult data={sr} />);
    });

    return (
      <div>
        {searchResults}
      </div>
    );
  }
});

var SearchResult = React.createClass({
  render: function () {
    var rawMarkup = this.props.data.preview;
    var genlink = function (permalink) {
      var target = searchTarget().split('%2F')[0];      
      if (target.startsWith('.')) {
        target = target.substring(4);
      } else if (target.endsWith('.%2B')) {
        target = target.substring(0, target.length - 4);
      }

      return '/read?v=' + permalink + '&hl=' + target;
    };

    return (
      <div className="search-result-item">
        <h4><a href={genlink(this.props.data.permalink)}>{this.props.data.title}</a></h4>
        <pre>
          <div className="search-result" dangerouslySetInnerHTML={{__html: rawMarkup}} />
        </pre>
      </div>
    );
  }
});

var LevelSwitchTab = React.createClass({
  render: function () {
    var count = this.props.levels;
    var active = this.props.activeLevel;
    var onSwitch = this.props.onSwitch;
    var nodes = function () {
      var results = new Array(count);
      for (var i = 1; i <= count; ++i) {
        var label = "Level " + i;
        if (i == active) {
          results[i-1] = (<li><span>{label}</span></li>);
        } else {
          results[i-1] = (<li><a href="javascript:;" onClick={onSwitch}>{label}</a></li>);
        }
      }
      return results;
    };

    return (
      <div>
        <div className="col-md-12 col-sm-12">
          <ul className="pagination pull-right">
            {nodes(count, active)}
          </ul>
        </div>
      </div>
    );
  }
});

var PaginationTab = React.createClass({
  prevPages: function () {
    if (this.props.startPage > 5) {
      this.props.setPages(this.props.startPage-5, 5, this.props.currentPage);
    }
  },

  nextPages: function () {
    this.props.setPages(this.props.startPage+5, 5, this.props.currentPage);
  },

  render: function() {
    var nodes = function () {
      var start = this.props.startPage;
      var n = this.props.numPages;
      var results = new Array(n);
      var current = this.props.currentPage;
      for (var i = 0; i < n; ++i) {
        if (i + start == current) {
          results[i] = (<li><span>{i + start}</span></li>);
        } else {
          results[i] = (<li><a href="javascript:;" onClick={this.props.gotoPage}>{i + start}</a></li>);
        }
      }
      return results;
    }.bind(this);

    return (
      <div className="row">
        <div className="col-md-4 col-sm-4 items-info">Items 1 to 9 of 10 total</div>
        <div className="col-md-8 col-sm-8">
          <ul className="pagination pull-right">
            <li><a href="javascript:;" onClick={this.prevPages}>«</a></li>
            {nodes()}
            <li><a href="javascript:;" onClick={this.nextPages}>»</a></li>
          </ul>
        </div>
      </div>
    );
  }
});

function searchByChar(char) {
  var url = apiHost + "/bychar?t=" + char + "&p=0&g=1" + "&" + apiKeyPostfix;
  React.render(
    <SearchView url={url} num_levels={3} />,
    document.getElementById('search-view')
  );
}

function searchByWord(word) {
  var url = apiHost + "/byvocabulary?t=" + word + "&p=0&g=1" + "&" + apiKeyPostfix;
  React.render(
    <SearchView url={url} num_levels={6} />,
    document.getElementById('search-view')
  );
}

function searchByGrammar(grammar) {
  var url = apiHost + "/bygrammar?t=" + grammar + "&p=0&g=1" + "&" + apiKeyPostfix;
  React.render(
    <SearchView url={url} num_levels={3} />,
    document.getElementById('search-view')
  );
}

function searchGeneric() {
  React.render(
    <SearchView />,
    document.getElementById('search-view')
  );
}


$(function () {
  var targetChar = getURLParameter("c");
  var targetVocabulary = getURLParameter("w");
  var targetGrammar = getURLParameter("m");

  if (typeof(targetChar) != "undefined") {
    searchByChar(targetChar);
  } else if (typeof(targetVocabulary) != "undefined") {
    searchByWord(targetVocabulary);
  } else if (typeof(targetGrammar) != "undefined") {
    searchByGrammar(targetGrammar);
  } else {
    searchGeneric();
  }
});


