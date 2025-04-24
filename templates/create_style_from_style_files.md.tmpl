---
page_title: "Provider: GeoServer - How to create style resource from style files"
description: |-
  A solution to generate geoserver_style resources from a folder with style file definitions
---

# The problem

On a project with a lot of styles, it can be tedious to explicitely declare a TF resource for each style you want to have in your GeoServer configuration.

# The solution

A solution to simplify the management of these resources can be the following:
- Use the `fileset` function to retrieve the list of the files in a reference folder
- Create a template `geoserver_style` resource with a `for_each` meta-argument to create the styles

Let's see an example on how to configure this.


```terraform
resource "geoserver_style" "automatic_styles" {
  # The list of the files we use to create our resources
  for_each = fileset(path.module,"styles/*.css")

  # Common elements
  workspace_name   = geoserver_workspace.nexsis.name
  format           = "css"
  version          = "1.0.0"

  # Specific elements from fileset
  name             = element(split(".",basename(each.key)),0)
  filename         = basename(each.key)
  style_definition = templatefile(format("%s/%s",path.module,each.key), { env = var.base_environment_name })

}
```

# Limit of the solution

`geoserver_style` resources created with this method cannot be referenced by other resources (for example in a group layer definition using non default styles). Following the ratio of referenced styles compared to the total number of styles, you can still use the previous pattern by removing the referenced styles from the `fileset` result, like this:

```terraform
resource "geoserver_style" "automatic_styles" {
  # The list of the files we use to create our resources, minus the ones we want to be able to reference
  for_each = setsubtract(fileset(path.module,"styles/*.css")
              ,["styles/limites_administratives.css","styles/ligne_electrique.css","styles/pylone.css"])

  # Common elements
  workspace_name   = geoserver_workspace.nexsis.name
  format           = "css"
  version          = "1.0.0"

  # Specific elements from fileset
  name             = element(split(".",basename(each.key)),0)
  filename         = basename(each.key)
  style_definition = templatefile(format("%s/%s",path.module,each.key), { env = var.base_environment_name })

}
```

# Similar use cases

This pattern could also be used for `geoserver_resource` resources.