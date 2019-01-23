package resourceReports

func GetHtmlTemplate() string {
	return `
		<!DOCTYPE html>
		<html>
		  <head>
		    <style>
		      table,
		      th,
		      td {
		        border: 1px solid black;
		      }
		    </style>
		  </head>
		  <body></body>
		  <script>
		    var ec2 = {{ if .EC2_JSON_PLACEHOLDER }} {{ .EC2_JSON_PLACEHOLDER }} {{ else }} undefined {{ end }}
		    var s3 = {{ if .S3_JSON_PLACEHOLDER }} {{ .S3_JSON_PLACEHOLDER }} {{ else }} undefined {{ end }}
		   
		    var HEADERS = {
		      volume_report: "Volume report",
		      instance_id: "Instance Id",
		      sortable_tags: "Sortable tags",
		      security_groups_ids: "Security groups ids",
		      availability_zone: "Availability zone",
		      name: "Bucket name",
		      encryption_type: "Default SSE",
		      logging_enabled: "Logging enabled",
		      acl_is_public: "ACL is public",
		      policy_is_public: "Policy is public"
		    };
		
		    function createHeader(name) {
		      document.body.appendChild(
		        document.createElement("h1").appendChild(document.createTextNode(name))
		      );
		    }
		
		    function createHeadRow(report, table) {
		      var tr = table.insertRow();
		      for (key in report[0]) {
		        var td = tr.insertCell();
		        td.appendChild(document.createTextNode(HEADERS[key]));
		      }
		    }
		
		    function createSubtableFromArray(array) {
		      var subtable = document.createElement("table");
		      array.forEach(function(value) {
		        var tr = subtable.insertRow();
		        var td = tr.insertCell();
		        td.appendChild(document.createTextNode(value));
		      });
		      return subtable;
		    }
		
		    function createSubtableFromMap(map) {
		      var subtable = document.createElement("table");
		      for (key in map) {
		        var tr = subtable.insertRow();
		        var tdkey = tr.insertCell();
		        tdkey.appendChild(document.createTextNode(key));
		        var tdvalue = tr.insertCell();
		        tdvalue.appendChild(document.createTextNode(map[key]));
		      }
		      return subtable;
		    }
		
		    function rowCreate(record, table) {
		      var tr = table.insertRow();
		      for (key in record) {
		        var td = tr.insertCell();
		        var value = record[key];
		        if (value !== undefined || value !== null) {
		          if (Array.isArray(value)) {
		            td.appendChild(createSubtableFromArray(value));
		          } else if (value instanceof Object) {
		            if (key === "sortable_tags") {
		              td.appendChild(createSubtableFromMap(value.Tags));
		            }
		          } else {
		            td.appendChild(document.createTextNode(record[key]));
		          }
		        } else {
		          td.appendChild(document.createTextNode(""));
		        }
		      }
		    }
		
		    function tableCreate(name, report) {
		      createHeader(name);
		      var table = document.createElement("table");
		      createHeadRow(report, table);
		      report.forEach(function(record) {
		        rowCreate(record, table);
		      });
		      document.body.appendChild(table);
		    }
		
		    var services = [{ name: "EC2", report: ec2 }, { name: "S3", report: s3 }];
		
		    services.forEach(function(service) {
		      if (service.name && service.report && Array.isArray(service.report)) {
		        tableCreate(service.name, service.report);
		      }
		    });
		  </script>
		</html>
	`
}
