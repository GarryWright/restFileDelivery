<?xml version="1.0" encoding="UTF-8"?>
<!-- This style sheet takes the Consolidated Instrument data within the document and places that data into the corresponding Template Variable fields withen the Template of the message -->
<xsl:transform version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:fn="http://www.w3.org/2005/02/xpath-functions"
    xmlns:xdt="http://www.w3.org/2005/02/xpath-datatypes" xmlns:xalan="http://xml.apache.org/xalan"
    exclude-result-prefixes="xalan">
    <xsl:output method="xml" omit-xml-declaration="yes" indent="yes"/>
    <xsl:template match="/*">
        <xsl:for-each select="*">
            <xsl:call-template name="getdesc"/>
        </xsl:for-each>
    </xsl:template>
    <xsl:template name="getdesc">
        <xsl:variable name="etag" select="name()"/>
        <xsl:element name="{$etag}">
            <xsl:for-each select="@*">
                <xsl:variable name="tag" select="name()"/>
                <xsl:attribute name="{$tag}">
                    <xsl:choose>
                        <xsl:when test="$tag='line-rate' and count(parent::*/descendant-or-self::*[@hits]) > 0">
                            <xsl:value-of
                                select="format-number(count(parent::*/descendant-or-self::*[@hits!='0']) div count(parent::*/descendant-or-self::*[@hits]) * 100,'###.##')"
                            />
                        </xsl:when>
                        <xsl:otherwise>
                            <xsl:value-of select="."/>
                        </xsl:otherwise>
                    </xsl:choose>                
                </xsl:attribute>         
            </xsl:for-each>
            <xsl:value-of select="text()"/>
            <xsl:for-each select="*">
                <xsl:call-template name="getdesc"/>
            </xsl:for-each>
            
        </xsl:element>
    </xsl:template>

</xsl:transform>
