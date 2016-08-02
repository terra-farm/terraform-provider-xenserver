/*
 * The MIT License (MIT)
 * Copyright (c) 2016 Maksym Borodin <borodin.maksym@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions
 * of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
 * THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 */
package main

import "github.com/hashicorp/terraform/helper/schema"

const (
	vdiSchemaUUID                 = "sr_uuid"
	vdiSchemaName                 = "name_label"
	vdiSchemaShared               = "shared"
	vdiSchemaRO                   = "read_only"
	vdiSchemaSize                 = "size"
)

func resourceVDI() *schema.Resource {
	return &schema.Resource{
		Create: resourceVDICreate,
		Read:   resourceVDIRead,
		Update: resourceVDIUpdate,
		Delete: resourceVDIDelete,
		Exists: resourceVDIExists,

		Schema: map[string]*schema.Schema{
			vdiSchemaUUID  : &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			vdiSchemaName  : &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			vdiSchemaShared: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},

			vdiSchemaRO    : &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},

			vdiSchemaSize  : &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

		},
	}
}

func resourceVDICreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)
}

func resourceVDIRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)
	return nil
}
func resourceVDIUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)
	return nil
}
func resourceVDIDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Connection)
	return nil
}
func resourceVDIExists(d *schema.ResourceData, m interface{}) (bool,error) {
	c := m.(*Connection)
	return false, nil
}