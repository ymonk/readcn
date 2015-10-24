/**
 * Created by myan on 21/09/15.
 */

var ArticleView = React.createClass({
    getInitialState: function () {
        return {data: {}};
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
        $("div.text-body").map(function () {
            var hl = decodeURI(getURLParameter('hl'));
            if (typeof(hl) == "undefined") {
                return;
            }
            var html = $(this).html();
            var re = new RegExp(hl, 'g');
            var newhtml = html.replace(re, '<span class="highlighted">' + hl + '</span>');
            $(this).html(newhtml);
        });
    },

    render: function () {
        var genlink = function (permalink) {
            return '/edit?v=' + permalink
        };

        if (typeof(this.state.data.body) == "undefined") {
            return (
                <div></div>
            );
        }
        var rawMarkup = this.state.data.body;
        return (
            <div>
                <div className="col-md-9 col-sm-9 blog-item">
                    <h4 className="active"><a href={genlink(this.state.data.permalink)}>{this.state.data.title}</a></h4>
                <pre>
                    <div className="text-body" dangerouslySetInnerHTML={{__html: rawMarkup}}/>
                </pre>
                    <ul className="blog-info">
                        <li><i className="fa fa-user"></i>{this.state.data.author}</li>
                        <li><i className="fa fa-calendar"></i>{this.state.data.publishedAt}</li>
                        <li><i className="fa fa-comments"></i>{this.state.data.num_comment}</li>
                        <li><i className="fa fa-tags"></i>{this.state.data.tags}</li>
                    </ul>
                </div>
            </div>
        );
    }
});

$(function () {
    // Get permalink from URL like http://localhost:5000/article?target=55e9423f5485093cbdfe835f
    var permalinkFromURL = function () {
        return getURLParameter("v")
    };

    var apiStartURL = function () {
        apiurl = apiHost + "/article/" + permalinkFromURL() + "?" + apiKeyPostfix;
        return apiurl;
    };
    React.render(
        <ArticleView url={apiStartURL()} />,
        document.getElementById('article-view')
    );
});


